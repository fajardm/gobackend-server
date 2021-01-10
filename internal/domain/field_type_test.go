package domain_test

import (
	"encoding/json"
	"testing"

	"github.com/fajardm/gobackend-server/internal/domain"
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
			var value domain.FieldType
			err := json.Unmarshal([]byte(tt.input), &value)

			actual := err.(errors.CustomError).Code()
			assert.Equal(t, tt.output, actual)
		})
	}
}

func TestFieldType_UnmarshalJSON_NoError(t *testing.T) {
	var flagtests = []struct {
		input  string
		output domain.FieldType
	}{
		{input: `"UUID"`, output: domain.FieldTypeUUID},
		{input: `"Boolean"`, output: domain.FieldTypeBoolean},
		{input: `"String"`, output: domain.FieldTypeString},
		{input: `"Decimal"`, output: domain.FieldTypeDecimal},
		{input: `"Integer"`, output: domain.FieldTypeInteger},
		{input: `"Date"`, output: domain.FieldTypeDate},
		{input: `"Array"`, output: domain.FieldTypeArray},
		{input: `"Object"`, output: domain.FieldTypeObject},
		{input: `"Pointer"`, output: domain.FieldTypePointer},
		{input: `"Relation"`, output: domain.FieldTypeRelation},
	}

	for _, tt := range flagtests {
		t.Run(tt.input, func(t *testing.T) {
			var value domain.FieldType
			err := json.Unmarshal([]byte(tt.input), &value)

			assert.NoError(t, err)
			assert.Equal(t, tt.output, value)
		})
	}
}

func TestFieldType_MarshalText_Error(t *testing.T) {
	var flagtests = []struct {
		name   string
		input  domain.FieldType
		output errors.Code
	}{
		{name: "0", input: domain.FieldType(0), output: errors.IncorrectFieldType},
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
		input  domain.FieldType
		output errors.Code
	}{
		{name: "uuid", input: domain.FieldTypeUUID, output: errors.IncorrectFieldType},
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
		input      domain.FieldType
		inputValue interface{}
		output     errors.Code
	}{
		{name: "string", input: domain.FieldTypeString, inputValue: 12, output: errors.IncorrectFieldType},
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
		input      domain.FieldType
		inputValue interface{}
		output     errors.Code
	}{
		{name: "string", input: domain.FieldTypeString, inputValue: "xxxx", output: errors.IncorrectFieldType},
	}

	for _, tt := range flagtests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.input.ValidateValue(tt.inputValue)

			assert.NoError(t, err)
		})
	}
}
