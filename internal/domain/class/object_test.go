package class_test

import (
	"encoding/json"
	"testing"

	"github.com/fajardm/gobackend-server/internal/domain/class"
	"github.com/fajardm/gobackend-server/internal/domain/schema"
	"github.com/fajardm/gobackend-server/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestOperator_UnmarshalJSON_Error(t *testing.T) {
	var flagtests = []struct {
		name   string
		input  string
		output errors.Code
	}{
		{
			name:   "test1",
			input:  `100`,
			output: errors.InvalidJSON,
		},
		{
			name:   "test2",
			input:  `"xxxx"`,
			output: errors.InvalidJSON,
		},
		{
			name:   "test3",
			input:  `{"username": "fajar"}`,
			output: errors.InvalidJSON,
		},
		{
			name:   "test4",
			input:  `{"address": 2}`,
			output: errors.IncorrectFieldType,
		},
	}

	for _, tt := range flagtests {
		t.Run(tt.name, func(t *testing.T) {
			schema := schema.Schema{
				Fields: schema.Fields{
					"address":   schema.Field{Type: schema.FieldTypeString, Required: false},
					"_rperm":    schema.Field{Type: schema.FieldTypeArray, Required: false},
					"_wperm":    schema.Field{Type: schema.FieldTypeArray, Required: false},
					"objectId":  schema.Field{Type: schema.FieldTypeUUID, Required: false},
					"createdAt": schema.Field{Type: schema.FieldTypeDate, Required: false},
					"updatedAt": schema.Field{Type: schema.FieldTypeDate, Required: false},
				},
				ClassName: "CustomUser",
			}
			object := class.NewObject(schema)
			err := json.Unmarshal([]byte(tt.input), &object)

			actual := err.(errors.CustomError).Code()
			assert.Equal(t, tt.output, actual)
		})
	}
}

func TestOperator_UnmarshalJSON_NoError(t *testing.T) {
	var flagtests = []struct {
		name   string
		input  string
		output class.ObjectValue
	}{
		{
			name: "test1",
			input: `
			{
				"address": "surakarta"
			}
			`,
			output: class.ObjectValue{
				schema.FieldName("address"): "surakarta",
			},
		},
	}

	for _, tt := range flagtests {
		t.Run(tt.name, func(t *testing.T) {
			schema := schema.Schema{
				Fields: schema.Fields{
					"address": schema.Field{Type: schema.FieldTypeString, Required: false},
				},
				ClassName: "CustomUser",
			}
			object := class.NewObject(schema)
			err := json.Unmarshal([]byte(tt.input), &object)

			assert.NoError(t, err)
			assert.Equal(t, tt.output, object.Value())
		})
	}
}

func TestOperator_MarshalJSON_NoError(t *testing.T) {
	var flagtests = []struct {
		name   string
		input  *class.Object
		output string
	}{
		{
			name: "test1",
			input: class.NewObjectWithValue(
				schema.Schema{
					Fields: schema.Fields{
						"address": schema.Field{Type: schema.FieldTypeString, Required: false},
					},
					ClassName: "CustomUser",
				},
				class.ObjectValue{
					schema.FieldName("address"): "surakarta",
				},
			),
			output: `{"address":"surakarta"}`,
		},
	}

	for _, tt := range flagtests {
		t.Run(tt.name, func(t *testing.T) {
			bytes, err := json.Marshal(tt.input)

			assert.NoError(t, err)
			assert.Equal(t, tt.output, string(bytes))
		})
	}
}
