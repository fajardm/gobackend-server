package domain_test

import (
	"encoding/json"
	"testing"

	"github.com/fajardm/gobackend-server/internal/domain"
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
			var value domain.ClassName
			err := json.Unmarshal([]byte(tt.input), &value)

			assert.Error(t, err)
		})
	}
}

func TestClassName_UnmarshalJSON_NoError(t *testing.T) {
	var flagtests = []struct {
		input  string
		output domain.ClassName
	}{
		{input: `"_User"`, output: domain.ClassName("_User")},
		{input: `"User"`, output: domain.ClassName("User")},
		{input: `"_Join:role:Student"`, output: domain.ClassName("_Join:role:Student")},
	}

	for _, tt := range flagtests {
		t.Run(tt.input, func(t *testing.T) {
			var value domain.ClassName
			err := json.Unmarshal([]byte(tt.input), &value)

			assert.NoError(t, err)
			assert.Equal(t, tt.output, value)
		})
	}
}

func TestJoinClassName(t *testing.T) {
	type _input struct {
		FieldName domain.FieldName
		ClassName domain.ClassName
	}

	var flagtests = []struct {
		name   string
		input  _input
		output domain.ClassName
	}{
		{
			name: "input1",
			input: _input{
				FieldName: domain.FieldName("role"),
				ClassName: domain.ClassName("Student"),
			},
			output: domain.ClassName("_Join:role:Student"),
		},
	}

	for _, tt := range flagtests {
		t.Run(tt.name, func(t *testing.T) {
			res := domain.JoinClassName(tt.input.FieldName, tt.input.ClassName)

			assert.Equal(t, tt.output, res)
		})
	}
}
