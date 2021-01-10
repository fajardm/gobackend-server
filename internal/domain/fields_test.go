package domain_test

import (
	"encoding/json"
	"testing"

	"github.com/fajardm/gobackend-server/internal/domain"
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
			var value domain.Field
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
		output domain.Field
	}{
		{
			name:  "input1",
			input: `{"type": "String", "required": true, "defaultValue": {"type": "String", "value": "xxxx"}}`,
			output: domain.Field{
				Type:     domain.FieldTypeString,
				Required: true,
				DefaultValue: domain.NewDefaultValue(domain.DefaultValueString{
					Type:  domain.FieldTypeString,
					Value: "xxxx",
				}),
			},
		},
	}

	for _, tt := range flagtests {
		t.Run(tt.name, func(t *testing.T) {
			var value domain.Field
			err := json.Unmarshal([]byte(tt.input), &value)

			assert.NoError(t, err)
			assert.Equal(t, tt.output, value)
		})
	}
}

func TestFields_Delete(t *testing.T) {
	fields := domain.Fields{
		domain.FieldName("xxxx"): domain.Field{
			Type:     domain.FieldTypeString,
			Required: true,
			DefaultValue: domain.NewDefaultValue(domain.DefaultValueString{
				Type:  domain.FieldTypeString,
				Value: "xxxx",
			}),
		},
	}

	_, ok := fields[domain.FieldName("xxxx")]
	assert.True(t, ok)

	fields.Delete(domain.FieldName("xxxx"))

	_, ok = fields[domain.FieldName("xxxx")]
	assert.False(t, ok)
}
