package domain_test

import (
	"encoding/json"
	"testing"

	"github.com/fajardm/gobackend-server/internal/domain"
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
			var value domain.FieldName
			err := json.Unmarshal([]byte(tt.input), &value)

			actual := err.(errors.CustomError).Code()
			assert.Equal(t, tt.output, actual)
		})
	}
}

func TestFieldName_UnmarshalJSON_NoError(t *testing.T) {
	var flagtests = []struct {
		input  string
		output domain.FieldName
	}{
		{input: `"Name"`, output: domain.FieldName("Name")},
		{input: `"name"`, output: domain.FieldName("name")},
		{input: `"name_"`, output: domain.FieldName("name_")},
	}

	for _, tt := range flagtests {
		t.Run(tt.input, func(t *testing.T) {
			var value domain.FieldName
			err := json.Unmarshal([]byte(tt.input), &value)

			assert.NoError(t, err)
			assert.Equal(t, tt.output, value)
		})
	}
}
