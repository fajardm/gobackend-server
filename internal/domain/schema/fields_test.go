package schema_test

import (
	"encoding/json"
	"testing"

	"github.com/fajardm/gobackend-server/internal/domain/schema"
	"github.com/fajardm/gobackend-server/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestField_UnmarshalJSON_Error(t *testing.T) {
	var flagtests = []struct {
		name   string
		input  string
		output errors.Code
	}{
		{
			name:   "input1",
			input:  `100`,
			output: errors.InvalidJSON,
		},
		{
			name:   "input2",
			input:  `{"type": "Pointer", "required": true}`,
			output: errors.MissingRequiredField,
		},
		{
			name:   "input3",
			input:  `{"type": "String", "required": true, "defaultValue": {"type": "Boolean"}}`,
			output: errors.IncorrectFieldType,
		},
		{
			name:   "input4",
			input:  `{"type": "Relation", "required": true, "targetClass": "User"}`,
			output: errors.IncorrectFieldType,
		},
	}

	for _, tt := range flagtests {
		t.Run(tt.name, func(t *testing.T) {
			var value schema.Field
			err := json.Unmarshal([]byte(tt.input), &value)

			actual := err.(errors.CustomError).Code()
			assert.Equal(t, tt.output, actual)
		})
	}
}

func TestField_UnmarshalJSON_NoError(t *testing.T) {
	var flagtests = []struct {
		name   string
		input  string
		output schema.Field
	}{
		{
			name:  "input1",
			input: `{"type": "String", "required": true, "defaultValue": {"type": "String", "value": "xxxx"}}`,
			output: schema.Field{
				Type:     schema.FieldTypeString,
				Required: true,
				DefaultValue: schema.NewDefaultValue(schema.DefaultValueString{
					Type:  schema.FieldTypeString,
					Value: "xxxx",
				}),
			},
		},
	}

	for _, tt := range flagtests {
		t.Run(tt.name, func(t *testing.T) {
			var value schema.Field
			err := json.Unmarshal([]byte(tt.input), &value)

			assert.NoError(t, err)
			assert.Equal(t, tt.output, value)
		})
	}
}

func TestFields_Delete(t *testing.T) {
	fields := schema.Fields{
		schema.FieldName("xxxx"): schema.Field{
			Type:     schema.FieldTypeString,
			Required: true,
			DefaultValue: schema.NewDefaultValue(schema.DefaultValueString{
				Type:  schema.FieldTypeString,
				Value: "xxxx",
			}),
		},
	}

	_, ok := fields[schema.FieldName("xxxx")]
	assert.True(t, ok)

	fields.Delete(schema.FieldName("xxxx"))

	_, ok = fields[schema.FieldName("xxxx")]
	assert.False(t, ok)
}
