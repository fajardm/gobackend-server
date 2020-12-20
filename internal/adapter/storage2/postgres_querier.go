package storage2

import "github.com/fajardm/gobackend-server/internal/domain/schema"

type PostgresQuerier interface {
	AllClasses() (string, []interface{}, error)
	ClassExists(className schema.ClassName) (string, []interface{}, error)
}

type postgresQuerier struct{}

func (postgresQuerier) AllClasses() (string, []interface{}, error) {
	query := `SELECT "className", "schema" FROM "_SCHEMA"`
	return query, nil, nil
}

func (postgresQuerier) ClassExists(className schema.ClassName) (string, []interface{}, error) {
	query := `SELECT exists (SELECT 1 FROM information_schema.tables WHERE table_name = $1)`
	args := []interface{}{className}
	return query, args, nil
}
