package postgres

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/fajardm/gobackend-server/internal/domain"
)

type Querier interface {
	ClassExist(className domain.ClassName) (string, []interface{}, error)
	GetAllClasses() (string, []interface{}, error)
	GetClass(className domain.ClassName) (string, []interface{}, error)
	CreateSchemaTableIfNotExists() (string, []interface{}, error)
	CreateClassTableIfNotExists(class domain.Class) (string, []interface{}, error)
	CreateRelationTableIfNotExists(className domain.ClassName, fieldName domain.FieldName) (string, []interface{}, error)
	CreateIndexIfNotExists(className domain.ClassName, indexName string, columns []domain.FieldName, unique bool) (string, []interface{}, error)
	CreateSchema(class domain.Class) (string, []interface{}, error)
	DropFieldIfExists(className domain.ClassName, fieldName domain.FieldName) (string, []interface{}, error)
	DropTablesIfExists(classNames ...domain.ClassName) (string, []interface{}, error)
}

var once sync.Once

var querierInstance Querier

func GetQuerier() Querier {
	if querierInstance == nil {
		once.Do(func() { querierInstance = new(querier) })
	}
	return querierInstance
}

type querier struct{}

func (querier) ClassExist(className domain.ClassName) (string, []interface{}, error) {
	query := `SELECT exists (SELECT 1 FROM information_schema.tables WHERE table_name = $1)`
	args := []interface{}{className}
	return query, args, nil
}

func (querier) GetAllClasses() (string, []interface{}, error) {
	query := `SELECT "className", "schema" FROM "_SCHEMA"`
	return query, nil, nil
}

func (querier) GetClass(className domain.ClassName) (string, []interface{}, error) {
	query := `SELECT "className", "schema" FROM "_SCHEMA" WHERE "className" = $1`
	return query, nil, nil
}

func (querier) CreateSchemaTableIfNotExists() (string, []interface{}, error) {
	query := `CREATE TABLE IF NOT EXISTS "_SCHEMA" ("className" varchar(120), "schema" jsonb PRIMARY KEY ("className"))`
	return query, nil, nil
}

func (querier) CreateClassTableIfNotExists(class domain.Class) (string, []interface{}, error) {
	query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS "%s"`, class.Name)
	columns := make([]string, 0)

	for fieldName, field := range class.Fields {
		if field.Type == domain.FieldTypeRelation {
			continue
		}

		columns = append(columns, fmt.Sprintf(`"%s"`, fieldName))

		pgType, err := NewFieldType(field.Type)
		if err != nil {
			return "", nil, err
		}
		columns = append(columns, pgType)

		if fieldName == domain.FieldObjectID {
			columns = append(columns, fmt.Sprintf(`PRIMARY KEY ("%s")`, fieldName))
		}
	}

	stringColumns := strings.Join(columns, ",")
	query = fmt.Sprintf("%s %s", query, stringColumns)

	return query, nil, nil
}

func (querier) CreateRelationTableIfNotExists(className domain.ClassName, fieldName domain.FieldName) (string, []interface{}, error) {
	query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS "%s" ("relatedId" uuid, "owningId" uuid, PRIMARY KEY("relatedId", "owningId"))`, domain.JoinClassName(fieldName, className))
	return query, nil, nil
}

func (querier) CreateIndexIfNotExists(className domain.ClassName, indexName string, columns []domain.FieldName, unique bool) (string, []interface{}, error) {
	sliceColumns := make([]string, 0)
	for _, column := range columns {
		sliceColumns = append(sliceColumns, column.String())
	}
	stringColumn := strings.Join(sliceColumns, ",")

	query := fmt.Sprintf(`CREATE INDEX IF NOT EXISTS "%s" ON "%s" ("%s")`, indexName, className, stringColumn)
	if unique {
		query = fmt.Sprintf(`%s %s`, query, "UNIQUE")
	}
	return query, nil, nil
}

func (querier) CreateSchema(class domain.Class) (string, []interface{}, error) {
	query := `INSERT INTO "_SCHEMA" ("className", "schema") VALUES ($1, $2)`

	jsonSchema, err := json.Marshal(class)
	if err != nil {
		return "", nil, err
	}

	args := []interface{}{class.Name, jsonSchema}
	return query, args, nil
}

func (querier) DropFieldIfExists(className domain.ClassName, fieldName domain.FieldName) (string, []interface{}, error) {
	query := fmt.Sprintf(`ALTER TABLE "%s" DROP COLUMN IF EXISTS "%s"`, className, fieldName)
	return query, nil, nil
}

func (querier) DropTablesIfExists(classNames ...domain.ClassName) (string, []interface{}, error) {
	sliceColumns := make([]string, 0)
	for _, className := range classNames {
		sliceColumns = append(sliceColumns, className.String())
	}
	stringClassNames := strings.Join(sliceColumns, ",")

	query := fmt.Sprintf(`DROP TABLE IF EXISTS "%s"`, stringClassNames)
	return query, nil, nil
}
