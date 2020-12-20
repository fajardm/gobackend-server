package schema_test

import (
	"encoding/json"
	"testing"

	"github.com/fajardm/gobackend-server/internal/domain/schema"
	"github.com/fajardm/gobackend-server/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestIndex_UnmarshalJSON_Error(t *testing.T) {
	var flagtests = []struct {
		name   string
		input  string
		output errors.Code
	}{
		{name: "input1", input: `{"columns": null, "unique": false}`, output: errors.InvalidJSON},
		{name: "input2", input: `{"columns": [], "unique": false}`, output: errors.InvalidJSON},
		{name: "input3", input: `100`, output: errors.InvalidJSON},
	}

	for _, tt := range flagtests {
		t.Run(tt.name, func(t *testing.T) {
			var value schema.Index
			err := json.Unmarshal([]byte(tt.input), &value)

			actual := err.(errors.CustomError).Code()
			assert.Equal(t, tt.output, actual)
		})
	}
}

func TestIndex_UnmarshalJSON_NoError(t *testing.T) {
	var flagtests = []struct {
		name   string
		input  string
		output schema.Index
	}{
		{name: "input1", input: `{"columns": ["name"], "unique": false}`, output: schema.Index{Columns: []schema.FieldName{"name"}, Unique: false}},
	}

	for _, tt := range flagtests {
		t.Run(tt.input, func(t *testing.T) {
			var value schema.Index
			err := json.Unmarshal([]byte(tt.input), &value)

			assert.NoError(t, err)
			assert.Equal(t, tt.output, value)
		})
	}
}

func TestIndexes_Delete(t *testing.T) {
	indexes := schema.Indexes{
		"xxxx": schema.Index{Columns: []schema.FieldName{"name"}, Unique: false},
	}

	assert.Equal(t, schema.FieldName("name"), indexes["xxxx"].Columns[0])

	indexes.Delete("xxxx")

	_, ok := indexes["xcxc"]
	assert.False(t, ok)
}
