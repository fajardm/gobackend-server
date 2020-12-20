package schema_test

import (
	"encoding/json"
	"testing"

	"github.com/fajardm/gobackend-server/internal/domain/schema"
	"github.com/stretchr/testify/assert"
)

func TestClassName_UnmarshalJSON_Error(t *testing.T) {
	var flagtests = []struct {
		input string
	}{
		{input: `"45User"`},
		{input: `"_Unknown"`},
		{input: `200`},
	}

	for _, tt := range flagtests {
		t.Run(tt.input, func(t *testing.T) {
			var value schema.ClassName
			err := json.Unmarshal([]byte(tt.input), &value)

			assert.Error(t, err)
		})
	}
}

func TestClassName_UnmarshalJSON_NoError(t *testing.T) {
	var flagtests = []struct {
		input  string
		output schema.ClassName
	}{
		{input: `"_User"`, output: schema.ClassName("_User")},
		{input: `"User"`, output: schema.ClassName("User")},
		{input: `"_Join:role:Student"`, output: schema.ClassName("_Join:role:Student")},
	}

	for _, tt := range flagtests {
		t.Run(tt.input, func(t *testing.T) {
			var value schema.ClassName
			err := json.Unmarshal([]byte(tt.input), &value)

			assert.NoError(t, err)
			assert.Equal(t, tt.output, value)
		})
	}
}

func TestJoinClassName(t *testing.T) {
	type _input struct {
		FieldName schema.FieldName
		ClassName schema.ClassName
	}

	var flagtests = []struct {
		name   string
		input  _input
		output schema.ClassName
	}{
		{
			name: "input1",
			input: _input{
				FieldName: schema.FieldName("role"),
				ClassName: schema.ClassName("Student"),
			},
			output: schema.ClassName("_Join:role:Student"),
		},
	}

	for _, tt := range flagtests {
		t.Run(tt.name, func(t *testing.T) {
			res := schema.JoinClassName(tt.input.FieldName, tt.input.ClassName)

			assert.Equal(t, tt.output, res)
		})
	}
}
