package schema_test

import (
	"encoding/json"
	"testing"

	"github.com/fajardm/gobackend-server/internal/domain/schema"
	"github.com/fajardm/gobackend-server/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestFieldName_UnmarshalJSON_Error(t *testing.T) {
	var flagtests = []struct {
		input  string
		output errors.Code
	}{
		{input: `"42_name"`, output: errors.InvalidFieldName},
		{input: `"_name"`, output: errors.InvalidFieldName},
		{input: `100`, output: errors.InvalidJSON},
	}

	for _, tt := range flagtests {
		t.Run(tt.input, func(t *testing.T) {
			var value schema.FieldName
			err := json.Unmarshal([]byte(tt.input), &value)

			actual := err.(errors.CustomError).Code()
			assert.Equal(t, tt.output, actual)
		})
	}
}

func TestFieldName_UnmarshalJSON_NoError(t *testing.T) {
	var flagtests = []struct {
		input  string
		output schema.FieldName
	}{
		{input: `"Name"`, output: schema.FieldName("Name")},
		{input: `"name"`, output: schema.FieldName("name")},
		{input: `"name_"`, output: schema.FieldName("name_")},
	}

	for _, tt := range flagtests {
		t.Run(tt.input, func(t *testing.T) {
			var value schema.FieldName
			err := json.Unmarshal([]byte(tt.input), &value)

			assert.NoError(t, err)
			assert.Equal(t, tt.output, value)
		})
	}
}
