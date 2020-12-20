package postgres

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/fajardm/gobackend-server/internal/domain/schema"
)

type SchemaQuerier interface {
	AllClasses() (string, []interface{}, error)
	FindClassByName(className schema.ClassName) (string, []interface{}, error)
	ExistsClass(className schema.ClassName) (string, []interface{}, error)
	CreateSchemaTableIfNotExists() (string, []interface{}, error)
	CreateClassTableIfNotExists(data schema.Schema) (string, []interface{}, error)
	CreateRelationTableIfNotExists(className schema.ClassName, fieldName schema.FieldName) (string, []interface{}, error)
	CreateIndexIfNotExists(className schema.ClassName, indexName string, columns []schema.FieldName, unique bool) (string, []interface{}, error)
	AddFieldIfNotExists(className schema.ClassName, fieldName schema.FieldName, field schema.Field) (string, []interface{}, error)
	CreateSchema(data schema.Schema) (string, []interface{}, error)
	UpdateSchema(data schema.Schema) (string, []interface{}, error)
	DropFieldIfExists(className schema.ClassName, fieldName schema.FieldName) (string, []interface{}, error)
	DropIndexesIfExists(indexes ...string) (string, []interface{}, error)
	DropTablesIfExists(classNames ...schema.ClassName) (string, []interface{}, error)
	DeleteSchema(className schema.ClassName) (string, []interface{}, error)
}

var once sync.Once

var schemaQuerierInstance SchemaQuerier

func GetSchemaQuerier() SchemaQuerier {
	if schemaQuerierInstance == nil {
		once.Do(func() { schemaQuerierInstance = new(schemaQuerier) })
	}
	return schemaQuerierInstance
}

type schemaQuerier struct{}

func (schemaQuerier) AllClasses() (string, []interface{}, error) {
	query := `SELECT "className", "schema" FROM "_SCHEMA"`
	return query, nil, nil
}

func (schemaQuerier) FindClassByName(className schema.ClassName) (string, []interface{}, error) {
	query := `SELECT "className", "schema" FROM "_SCHEMA" WHERE "className" = $1`
	return query, nil, nil
}

func (schemaQuerier) ExistsClass(className schema.ClassName) (string, []interface{}, error) {
	query := `SELECT exists (SELECT 1 FROM information_schema.tables WHERE table_name = $1)`
	args := []interface{}{className}
	return query, args, nil
}

func (schemaQuerier) CreateSchemaTableIfNotExists() (string, []interface{}, error) {
	query := `CREATE TABLE IF NOT EXISTS "_SCHEMA" ("className" varchar(120), "schema" jsonb PRIMARY KEY ("className"))`
	return query, nil, nil
}

func (schemaQuerier) CreateClassTableIfNotExists(data schema.Schema) (string, []interface{}, error) {
	query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS "%s"`, data.ClassName)
	columns := make([]string, 0)

	for fieldName, field := range data.Fields {
		if field.Type == schema.FieldTypeRelation {
			continue
		}

		columns = append(columns, fmt.Sprintf(`"%s"`, fieldName))

		pgType, err := NewFieldType(field.Type)
		if err != nil {
			return "", nil, err
		}
		columns = append(columns, pgType)

		if fieldName == schema.FieldObjectID {
			columns = append(columns, fmt.Sprintf(`PRIMARY KEY ("%s")`, fieldName))
		}
	}

	stringColumns := strings.Join(columns, ",")
	query = fmt.Sprintf("%s %s", query, stringColumns)

	return query, nil, nil
}

func (schemaQuerier) CreateRelationTableIfNotExists(className schema.ClassName, fieldName schema.FieldName) (string, []interface{}, error) {
	query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS "_Join:%s:%s" ("relatedId" uuid, "owningId" uuid, PRIMARY KEY("relatedId", "owningId"))`, fieldName, className)
	return query, nil, nil
}

func (schemaQuerier) CreateIndexIfNotExists(className schema.ClassName, indexName string, columns []schema.FieldName, unique bool) (string, []interface{}, error) {
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

func (schemaQuerier) AddFieldIfNotExists(className schema.ClassName, fieldName schema.FieldName, field schema.Field) (string, []interface{}, error) {
	pgType, err := NewFieldType(field.Type)
	if err != nil {
		return "", nil, err
	}
	query := fmt.Sprintf(`ALTER TABLE "%s" ADD COLUMN IF NOT EXISTS "%s" %s`, className, fieldName, pgType)
	return query, nil, nil
}

func (schemaQuerier) CreateSchema(data schema.Schema) (string, []interface{}, error) {
	query := `INSERT INTO "_SCHEMA" ("className", "schema") VALUES ($1, $2)`

	jsonSchema, err := json.Marshal(data)
	if err != nil {
		return "", nil, err
	}

	args := []interface{}{data.ClassName, jsonSchema}
	return query, args, nil
}

func (schemaQuerier) UpdateSchema(data schema.Schema) (string, []interface{}, error) {
	query := `UPDATE "%s" SET "schema" = $1 WHERE "className" = $2`

	jsonSchema, err := json.Marshal(data)
	if err != nil {
		return "", nil, err
	}

	args := []interface{}{jsonSchema, data.ClassName}
	return query, args, nil
}

func (schemaQuerier) DropFieldIfExists(className schema.ClassName, fieldName schema.FieldName) (string, []interface{}, error) {
	query := fmt.Sprintf(`ALTER TABLE "%s" DROP COLUMN IF EXISTS "%s"`, className, fieldName)
	return query, nil, nil
}

func (schemaQuerier) DropIndexesIfExists(indexes ...string) (string, []interface{}, error) {
	query := fmt.Sprintf(`DROP INDEX IF EXISTS "%s"`, strings.Join(indexes, `","`))
	return query, nil, nil
}

func (schemaQuerier) DropTablesIfExists(classNames ...schema.ClassName) (string, []interface{}, error) {
	sliceColumns := make([]string, 0)
	for _, className := range classNames {
		sliceColumns = append(sliceColumns, className.String())
	}
	stringClassNames := strings.Join(sliceColumns, ",")

	query := fmt.Sprintf(`DROP TABLE IF EXISTS "%s"`, stringClassNames)
	return query, nil, nil
}

func (schemaQuerier) DeleteSchema(className schema.ClassName) (string, []interface{}, error) {
	query := fmt.Sprintf(`DELETE FROM "_SCHEMA" WHERE "className" = $1`)
	args := []interface{}{className}
	return query, args, nil
}
