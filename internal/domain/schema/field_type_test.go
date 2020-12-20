package schema_test

import (
	"encoding/json"
	"testing"

	"github.com/fajardm/gobackend-server/internal/domain/schema"
	"github.com/fajardm/gobackend-server/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestFieldType_UnmarshalJSON_Error(t *testing.T) {
	var flagtests = []struct {
		input  string
		output errors.Code
	}{
		{input: `"Int64"`, output: errors.IncorrectFieldType},
		{input: `100`, output: errors.InvalidJSON},
	}

	for _, tt := range flagtests {
		t.Run(tt.input, func(t *testing.T) {
			var value schema.FieldType
			err := json.Unmarshal([]byte(tt.input), &value)

			actual := err.(errors.CustomError).Code()
			assert.Equal(t, tt.output, actual)
		})
	}
}

func TestFieldType_UnmarshalJSON_NoError(t *testing.T) {
	var flagtests = []struct {
		input  string
		output schema.FieldType
	}{
		{input: `"UUID"`, output: schema.FieldTypeUUID},
		{input: `"Boolean"`, output: schema.FieldTypeBoolean},
		{input: `"String"`, output: schema.FieldTypeString},
		{input: `"Decimal"`, output: schema.FieldTypeDecimal},
		{input: `"Integer"`, output: schema.FieldTypeInteger},
		{input: `"Date"`, output: schema.FieldTypeDate},
		{input: `"Array"`, output: schema.FieldTypeArray},
		{input: `"Object"`, output: schema.FieldTypeObject},
		{input: `"Pointer"`, output: schema.FieldTypePointer},
		{input: `"Relation"`, output: schema.FieldTypeRelation},
	}

	for _, tt := range flagtests {
		t.Run(tt.input, func(t *testing.T) {
			var value schema.FieldType
			err := json.Unmarshal([]byte(tt.input), &value)

			assert.NoError(t, err)
			assert.Equal(t, tt.output, value)
		})
	}
}

func TestFieldType_MarshalText_Error(t *testing.T) {
	var flagtests = []struct {
		name   string
		input  schema.FieldType
		output errors.Code
	}{
		{name: "0", input: schema.FieldType(0), output: errors.IncorrectFieldType},
	}

	for _, tt := range flagtests {
		t.Run(tt.name, func(t *testing.T) {
			bytes, err := json.Marshal(tt.input)

			assert.Error(t, err)
			assert.Equal(t, []byte(nil), bytes)
		})
	}
}

func TestFieldType_MarshalText_NoError(t *testing.T) {
	var flagtests = []struct {
		name   string
		input  schema.FieldType
		output errors.Code
	}{
		{name: "uuid", input: schema.FieldTypeUUID, output: errors.IncorrectFieldType},
	}

	for _, tt := range flagtests {
		t.Run(tt.name, func(t *testing.T) {
			bytes, err := json.Marshal(tt.input)

			assert.NoError(t, err)
			assert.Equal(t, `"UUID"`, string(bytes))
		})
	}
}

func TestFieldType_ValidateValue_Error(t *testing.T) {
	var flagtests = []struct {
		name       string
		input      schema.FieldType
		inputValue interface{}
		output     errors.Code
	}{
		{name: "string", input: schema.FieldTypeString, inputValue: 12, output: errors.IncorrectFieldType},
	}

	for _, tt := range flagtests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.input.ValidateValue(tt.inputValue)

			actual := err.(errors.CustomError).Code()
			assert.Equal(t, tt.output, actual)
		})
	}
}

func TestFieldType_ValidateValue_NoError(t *testing.T) {
	var flagtests = []struct {
		name       string
		input      schema.FieldType
		inputValue interface{}
		output     errors.Code
	}{
		{name: "string", input: schema.FieldTypeString, inputValue: "xxxx", output: errors.IncorrectFieldType},
	}

	for _, tt := range flagtests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.input.ValidateValue(tt.inputValue)

			assert.NoError(t, err)
		})
	}
}
