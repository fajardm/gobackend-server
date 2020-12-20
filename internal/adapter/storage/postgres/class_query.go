package postgres

import (
	"fmt"
	"strings"

	"github.com/fajardm/gobackend-server/internal/domain/class"
	"github.com/fajardm/gobackend-server/internal/domain/schema"
)

type ClassQuery interface {
	BuildWhereClause(schema schema.Schema, query interface{}, index int, caseInsensitive bool)
}

type classQuery struct{}

func (classQuery) createObject(sch schema.Schema, object class.Object) (string, []interface{}, error) {
	columns, args := make([]string, 0), make([]interface{}, 0)

	for fieldName, value := range object.Value() {
		columns, args = append(columns, fieldName.String()), append(args, value)
	}

	stringColumns, stringValues := fmt.Sprintf(`"%s"`, strings.Join(columns, `","`)), ""

	for i, _ := range args {
		stringValues = stringValues + fmt.Sprintf("$%d", i)
	}

	query := fmt.Sprintf(`INSERT INTO "%s" (%s) VALUES (%s)`, sch.ClassName, stringColumns, stringValues)
	return query, args, nil
}

// func (q classQuery) BuildWhereClause(schema schema.Schema, query class.Query, index int, caseInsensitive bool) {
// 	patterns, values, sorts := make([]string, 0), make([]string, 0), make([]string, 0)

// 	for field, comparator := range query {

// 	}
// }
