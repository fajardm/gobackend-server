package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/fajardm/gobackend-server/internal/adapter/database"
	"github.com/fajardm/gobackend-server/internal/domain/schema"
	"github.com/fajardm/gobackend-server/pkg/errors"
	"github.com/fajardm/gobackend-server/pkg/trx"
	"github.com/jmoiron/sqlx"
)

var _ schema.Repository = (*Schema)(nil)

type Schema struct {
	DB        *database.Postgres `inject:"postgres"`
	Unmarshal Unmarshal          `inject:"unmarshal"`
	Querier   SchemaQuerier      `inject:"schemaQuerier"`
}

func (r Schema) All(ctx context.Context) (schema.Schemas, error) {
	query, args, err := r.Querier.AllClasses()
	if err != nil {
		return nil, err
	}

	rows, err := r.DB.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make(schema.Schemas)
	for rows.Next() {
		var className schema.ClassName
		var schemaByte []byte
		if err := rows.Scan(&className, &schemaByte); err != nil {
			return nil, err
		}

		var schema schema.Schema
		if err := r.Unmarshal(schemaByte, &schema); err != nil {
			return nil, err
		}
		res[className] = schema
	}

	return res, nil
}

func (r Schema) FindByClassName(ctx context.Context, className schema.ClassName) (*schema.Schema, error) {
	query, args, err := r.Querier.FindClassByName(className)
	if err != nil {
		return nil, err
	}

	var schemaByte []byte
	if err := r.DB.QueryRowxContext(ctx, query, args...).Scan(&className, &schemaByte); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New(errors.DataNotFound, fmt.Sprintf("class %s not exists", className))
		}
		return nil, err
	}

	var res *schema.Schema
	if err := r.Unmarshal(schemaByte, &res); err != nil {
		return nil, err
	}

	return res, nil
}

func (r Schema) Exists(ctx context.Context, className schema.ClassName) (bool, error) {
	query, args, err := r.Querier.ExistsClass(className)
	if err != nil {
		return false, err
	}

	var res bool
	if err := r.DB.QueryRowxContext(ctx, query, args...).Scan(&res); err != nil {
		return false, err
	}

	return res, nil
}

func (r Schema) Create(ctx context.Context, data schema.Schema) error {
	err := trx.WithTx(r.DB.DB, func(tx *sqlx.Tx, errChan chan error) {
		if err := r.txCreateSchemaTableIfNotExists(ctx, tx); err != nil {
			errChan <- err
			return
		}

		if err := r.txCreateClassTableIfNotExists(ctx, tx, data); err != nil {
			errChan <- err
			return
		}

		if err := r.txCreateRelationTablesIfNotExists(ctx, tx, data); err != nil {
			errChan <- err
			return
		}

		if err := r.txCreateIndexesIfNotExists(ctx, tx, data.ClassName, data.Indexes); err != nil {
			errChan <- err
			return
		}

		if err := r.txCreateSchema(ctx, tx, data); err != nil {
			errChan <- err
			return
		}
	})
	return err
}

func (r Schema) Update(ctx context.Context, data schema.Schema) error {
	existing, err := r.FindByClassName(ctx, data.ClassName)
	if err != nil {
		return err
	}

	errChan := make(chan error)

	type DoField struct {
		key   schema.FieldName
		field schema.Field
	}

	insertFieldChan := make(chan DoField)

	go func() {
		for key, field := range data.Fields {
			if _, ok := schema.DefaultColumn[key]; ok {
				errChan <- errors.New(errors.ChangedImmutableField, fmt.Sprintf("field %s cannot be changed", key))
				return
			}

			if _, ok := schema.DefaultColumns[data.ClassName.String()][key]; ok {
				errChan <- errors.New(errors.ChangedImmutableField, fmt.Sprintf("field %s cannot be changed", key))
				return
			}

			if _, ok := existing.Fields[key]; ok {
				errChan <- errors.New(errors.ChangedImmutableField, fmt.Sprintf("field %s cannot be changed", key))
				return
			}

			if _, ok := existing.Fields[key]; !ok {
				insertFieldChan <- DoField{key: key, field: field}
			}
		}
	}()

	deleteFieldChan := make(chan DoField)

	go func() {
		for key, field := range existing.Fields {
			if _, ok := data.Fields[key]; !ok {
				deleteFieldChan <- DoField{key: key, field: field}
			}
		}
		close(deleteFieldChan)
	}()

	type DoIndex struct {
		key   string
		index schema.Index
	}

	insertIndexChan := make(chan DoIndex)

	go func() {
		for key, index := range data.Indexes {
			if _, ok := existing.Indexes[key]; !ok {
				insertIndexChan <- DoIndex{key: key, index: index}
			}
		}
	}()

	deleteIndexChan := make(chan DoIndex)

	go func() {
		for key, index := range existing.Indexes {
			if _, ok := data.Indexes[key]; !ok {
				deleteIndexChan <- DoIndex{key: key, index: index}
			}
		}
		close(deleteIndexChan)
	}()

	deleteFields := make(schema.Fields)
	insertFields := make(schema.Fields)

	deleteIndexes := make([]string, 0)
	toInsertIndexes := make(schema.Indexes)

	for {
		select {
		case err := <-errChan:
			return err
		case deleteField := <-deleteFieldChan:
			deleteFields[deleteField.key] = deleteField.field
			existing.Fields.Delete(deleteField.key)
		case deleteIndex := <-deleteIndexChan:
			deleteIndexes = append(deleteIndexes, deleteIndex.key)
			existing.Indexes.Delete(deleteIndex.key)
		case inserField := <-insertFieldChan:
			insertFields[inserField.key] = inserField.field
			existing.Fields[inserField.key] = inserField.field
		case insertIndex := <-insertIndexChan:
			toInsertIndexes[insertIndex.key] = insertIndex.index
			existing.Indexes[insertIndex.key] = insertIndex.index
		}

		if deleteFieldChan == nil && deleteIndexChan == nil {
			break
		}
	}

	err = trx.WithTx(r.DB.DB, func(tx *sqlx.Tx, errChan chan error) {
		if err := r.txDropFieldsAndRelations(ctx, tx, data.ClassName, deleteFields); err != nil {
			errChan <- err
			return
		}

		if err := r.txAddFieldsAndRelationsIfNotExists(ctx, tx, data.ClassName, insertFields); err != nil {
			errChan <- err
			return
		}

		if err := r.txDropIndexesIfExists(ctx, tx, deleteIndexes...); err != nil {
			errChan <- err
			return
		}

		if err := r.txCreateIndexesIfNotExists(ctx, tx, data.ClassName, toInsertIndexes); err != nil {
			errChan <- err
			return
		}

		if err := r.txUpdateSchema(ctx, tx, *existing); err != nil {
			errChan <- err
			return
		}
	})
	return err
}

func (r Schema) Delete(ctx context.Context, className schema.ClassName) error {
	existing, err := r.FindByClassName(ctx, className)
	if err != nil {
		return err
	}

	err = trx.WithTx(r.DB.DB, func(tx *sqlx.Tx, errChan chan error) {
		if err := r.txDropTables(ctx, tx, className); err != nil {
			errChan <- err
			return
		}

		for key, field := range existing.Fields {
			if field.Type == schema.FieldTypeRelation {
				if err = r.txDropTables(ctx, tx, schema.JoinClassName(key, *field.TargetClass)); err != nil {
					errChan <- err
					return
				}
			}
		}

		if err := r.txDeleteSchema(ctx, tx, className); err != nil {
			errChan <- err
			return
		}
	})
	return err
}

func (r Schema) txCreateSchemaTableIfNotExists(ctx context.Context, tx *sqlx.Tx) error {
	query, agrs, err := r.Querier.CreateSchemaTableIfNotExists()
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, query, agrs...)
	return err
}

func (r Schema) txCreateClassTableIfNotExists(ctx context.Context, tx *sqlx.Tx, data schema.Schema) error {
	query, args, err := r.Querier.CreateClassTableIfNotExists(data)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, query, args...)
	return err
}

func (r Schema) txCreateRelationTablesIfNotExists(ctx context.Context, tx *sqlx.Tx, data schema.Schema) (err error) {
	for fieldName, field := range data.Fields {
		if field.Type == schema.FieldTypeRelation {
			if err = r.txCreateRelationTableIfNotExists(ctx, tx, data.ClassName, fieldName); err != nil {
				return err
			}
		}
	}
	return
}

func (r Schema) txCreateRelationTableIfNotExists(ctx context.Context, tx *sqlx.Tx, className schema.ClassName, fieldName schema.FieldName) error {
	query, args, err := r.Querier.CreateRelationTableIfNotExists(className, fieldName)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, query, args...)
	return err
}

func (r Schema) txCreateIndexesIfNotExists(ctx context.Context, tx *sqlx.Tx, className schema.ClassName, indexes schema.Indexes) error {
	for indexName, index := range indexes {
		if err := r.txCreateIndexIfNotExists(ctx, tx, className, indexName, index); err != nil {
			return err
		}
	}
	return nil
}

func (r Schema) txAddFieldIfNotExists(ctx context.Context, tx *sqlx.Tx, className schema.ClassName, fieldName schema.FieldName, field schema.Field) error {
	query, args, err := r.Querier.AddFieldIfNotExists(className, fieldName, field)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, query, args...)
	return err
}

func (r Schema) txAddFieldsAndRelationsIfNotExists(ctx context.Context, tx *sqlx.Tx, className schema.ClassName, fields schema.Fields) error {
	for name, field := range fields {
		if field.Type != schema.FieldTypeRelation {
			if err := r.txAddFieldIfNotExists(ctx, tx, className, name, field); err != nil {
				return err
			}
		} else {
			if err := r.txCreateRelationTableIfNotExists(ctx, tx, className, name); err != nil {
				return err
			}
		}
	}
	return nil
}

func (r Schema) txCreateIndexIfNotExists(ctx context.Context, tx *sqlx.Tx, className schema.ClassName, indexName string, index schema.Index) error {
	query, args, err := r.Querier.CreateIndexIfNotExists(className, indexName, index.Columns, index.Unique)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, query, args...)
	return err
}

func (r Schema) txCreateSchema(ctx context.Context, tx *sqlx.Tx, data schema.Schema) error {
	query, args, err := r.Querier.CreateSchema(data)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, query, args...)
	return err
}

func (r Schema) txUpdateSchema(ctx context.Context, tx *sqlx.Tx, data schema.Schema) error {
	query, args, err := r.Querier.UpdateSchema(data)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, query, args...)
	return err
}

func (r Schema) txDropField(ctx context.Context, tx *sqlx.Tx, className schema.ClassName, fieldName schema.FieldName) error {
	query, agrs, err := r.Querier.DropFieldIfExists(className, fieldName)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, query, agrs...)
	return err
}

func (r Schema) txDropFieldsAndRelations(ctx context.Context, tx *sqlx.Tx, className schema.ClassName, fields schema.Fields) error {
	for key, field := range fields {
		if field.Type != schema.FieldTypeRelation {
			if err := r.txDropField(ctx, tx, className, key); err != nil {
				return err
			}
		} else {
			if err := r.txDropTables(ctx, tx, schema.JoinClassName(key, className)); err != nil {
				return err
			}
		}
	}
	return nil
}

func (r Schema) txDropIndexesIfExists(ctx context.Context, tx *sqlx.Tx, indexes ...string) error {
	query, agrs, err := r.Querier.DropIndexesIfExists(indexes...)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, query, agrs...)
	return err
}

func (r Schema) txDropTables(ctx context.Context, tx *sqlx.Tx, classNames ...schema.ClassName) error {
	query, agrs, err := r.Querier.DropTablesIfExists(classNames...)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, query, agrs...)
	return err
}

func (r Schema) txDeleteSchema(ctx context.Context, tx *sqlx.Tx, className schema.ClassName) error {
	query, agrs, err := r.Querier.DeleteSchema(className)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, query, agrs...)
	return err
}
