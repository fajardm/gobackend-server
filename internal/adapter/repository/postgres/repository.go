package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/fajardm/gobackend-server/internal/domain"
	"github.com/fajardm/gobackend-server/pkg/errors"
	"github.com/fajardm/gobackend-server/pkg/trx"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	DB        *sqlx.DB  `inject:"sqlx"`
	Unmarshal Unmarshal `inject:"unmarshal"`
	Querier   Querier   `inject:"querier"`
}

func (r Repository) ClassExist(ctx context.Context, className domain.ClassName) (bool, error) {
	query, args, err := r.Querier.ClassExist(className)
	if err != nil {
		return false, err
	}

	var res bool
	if err := r.DB.QueryRowxContext(ctx, query, args...).Scan(&res); err != nil {
		return false, err
	}

	return res, nil
}

func (r Repository) GetAllClasses(ctx context.Context) (domain.Classes, error) {
	query, args, err := r.Querier.GetAllClasses()
	if err != nil {
		return nil, err
	}

	rows, err := r.DB.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make(domain.Classes)
	for rows.Next() {
		var (
			className domain.ClassName
			classByte []byte
			class     domain.Class
		)

		if err := rows.Scan(&className, &classByte); err != nil {
			return nil, err
		}

		if err := r.Unmarshal(classByte, &class); err != nil {
			return nil, err
		}
		res[className] = class
	}

	return res, nil
}

func (r Repository) GetClass(ctx context.Context, className domain.ClassName) (*domain.Class, error) {
	query, args, err := r.Querier.GetClass(className)
	if err != nil {
		return nil, err
	}

	var (
		classByte []byte
		res       *domain.Class
	)

	if err := r.DB.QueryRowxContext(ctx, query, args...).Scan(&className, &classByte); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New(errors.DataNotFound, fmt.Sprintf("class %s not exists", className))
		}
		return nil, err
	}

	if err := r.Unmarshal(classByte, &res); err != nil {
		return nil, err
	}

	return res, nil
}

func (r Repository) CreateClass(ctx context.Context, class domain.Class) error {
	err := trx.WithTx(r.DB, func(tx *sqlx.Tx, errChan chan error) {
		ctx = context.WithValue(ctx, "tx", tx)

		if err := r.createClassTableIfNotExists(ctx, class); err != nil {
			errChan <- err
			return
		}

		if err := r.createRelationTablesIfNotExists(ctx, class); err != nil {
			errChan <- err
			return
		}

		if err := r.createIndexesIfNotExists(ctx, class.Name, class.Indexes); err != nil {
			errChan <- err
			return
		}

		if err := r.createSchema(ctx, class); err != nil {
			errChan <- err
			return
		}
	})
	return err
}

func (r Repository) UpdateClass(ctx context.Context, class domain.Class) error {
	err := trx.WithTx(r.DB, func(tx *sqlx.Tx, errChan chan error) {
		if err := r.dropFieldsAndRelations(ctx, class.Name, class.ToUpdate.DeleteFields); err != nil {
			errChan <- err
			return
		}
	})
	return err
}

func (r Repository) createSchemaTableIfNotExists(ctx context.Context) error {
	query, agrs, err := r.Querier.CreateSchemaTableIfNotExists()
	if err != nil {
		return err
	}
	_, err = r.DB.ExecContext(ctx, query, agrs...)
	return err
}

func (r Repository) createClassTableIfNotExists(ctx context.Context, class domain.Class) error {
	query, args, err := r.Querier.CreateClassTableIfNotExists(class)
	if err != nil {
		return err
	}
	execer := r.getExecer(ctx)
	_, err = execer.ExecContext(ctx, query, args...)
	return err
}

func (r Repository) createRelationTablesIfNotExists(ctx context.Context, class domain.Class) (err error) {
	for fieldName, field := range class.Fields {
		if field.Type == domain.FieldTypeRelation {
			if err = r.createRelationTableIfNotExists(ctx, class.Name, fieldName); err != nil {
				return err
			}
		}
	}
	return
}

func (r Repository) createRelationTableIfNotExists(ctx context.Context, className domain.ClassName, fieldName domain.FieldName) error {
	query, args, err := r.Querier.CreateRelationTableIfNotExists(className, fieldName)
	if err != nil {
		return err
	}
	execer := r.getExecer(ctx)
	_, err = execer.ExecContext(ctx, query, args...)
	return err
}

func (r Repository) createIndexesIfNotExists(ctx context.Context, className domain.ClassName, indexes domain.Indexes) error {
	for indexName, index := range indexes {
		if err := r.createIndexIfNotExists(ctx, className, indexName, index); err != nil {
			return err
		}
	}
	return nil
}

func (r Repository) createIndexIfNotExists(ctx context.Context, className domain.ClassName, indexName string, index domain.Index) error {
	query, args, err := r.Querier.CreateIndexIfNotExists(className, indexName, index.Columns, index.Unique)
	if err != nil {
		return err
	}
	execer := r.getExecer(ctx)
	_, err = execer.ExecContext(ctx, query, args...)
	return err
}

func (r Repository) createSchema(ctx context.Context, class domain.Class) error {
	query, args, err := r.Querier.CreateSchema(class)
	if err != nil {
		return err
	}
	execer := r.getExecer(ctx)
	_, err = execer.ExecContext(ctx, query, args...)
	return err
}

func (r Repository) dropFieldsAndRelations(ctx context.Context, className domain.ClassName, fields domain.Fields) error {
	for key, field := range fields {
		if field.Type != domain.FieldTypeRelation {
			if err := r.dropField(ctx, className, key); err != nil {
				return err
			}
		} else {
			if err := r.dropTables(ctx, domain.JoinClassName(key, className)); err != nil {
				return err
			}
		}
	}
	return nil
}

func (r Repository) dropField(ctx context.Context, className domain.ClassName, fieldName domain.FieldName) error {
	query, agrs, err := r.Querier.DropFieldIfExists(className, fieldName)
	if err != nil {
		return err
	}
	execer := r.getExecer(ctx)
	_, err = execer.ExecContext(ctx, query, agrs...)
	return err
}

func (r Repository) dropTables(ctx context.Context, classNames ...domain.ClassName) error {
	query, agrs, err := r.Querier.DropTablesIfExists(classNames...)
	if err != nil {
		return err
	}
	execer := r.getExecer(ctx)
	_, err = execer.ExecContext(ctx, query, agrs...)
	return err
}

func (r Repository) getExecer(ctx context.Context) sqlx.ExecerContext {
	txInterface := ctx.Value("tx")
	tx, ok := txInterface.(*sqlx.Tx)
	if !ok {
		return r.DB
	}
	return tx
}
