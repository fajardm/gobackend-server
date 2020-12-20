package schema_test

import (
	"encoding/json"
	"testing"

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
			name:   "input1",
			input:  `100`,
			output: errors.InvalidJSON,
		},
		{
			name:   "input2",
			input:  `"xxx"`,
			output: errors.IncorrectOperation,
		},
	}

	for _, tt := range flagtests {
		t.Run(tt.name, func(t *testing.T) {
			var value schema.Operator
			err := json.Unmarshal([]byte(tt.input), &value)

			actual := err.(errors.CustomError).Code()
			assert.Equal(t, tt.output, actual)
		})
	}
}

func TestOperator_UnmarshalJSON_NoError(t *testing.T) {
	var flagtests = []struct {
		name   string
		input  string
		output schema.Operator
	}{
		{
			name:   "input1",
			input:  `"delete"`,
			output: schema.OperatorDelete,
		},
	}

	for _, tt := range flagtests {
		t.Run(tt.name, func(t *testing.T) {
			var value schema.Operator
			err := json.Unmarshal([]byte(tt.input), &value)

			assert.NoError(t, err)
			assert.Equal(t, tt.output, value)
		})
	}
}

func TestOperator_MarshalText_Error(t *testing.T) {
	var flagtests = []struct {
		name  string
		input schema.Operator
	}{
		{
			name:  "input1",
			input: schema.Operator(0),
		},
	}

	for _, tt := range flagtests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := json.Marshal(tt.input)

			assert.Error(t, err)
			assert.Nil(t, res)
		})
	}
}

func TestOperator_MarshalJSON_Error(t *testing.T) {
	type Test struct {
		Operator schema.Operator `json:"operator"`
	}

	var flagtests = []struct {
		name  string
		input Test
	}{
		{
			name:  "input1",
			input: Test{Operator: schema.Operator(0)},
		},
	}

	for _, tt := range flagtests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := json.Marshal(tt.input)

			assert.Error(t, err)
			assert.Nil(t, res)
		})
	}
}

func TestOperator_MarshalText_NoError(t *testing.T) {
	var flagtests = []struct {
		name   string
		input  schema.Operator
		output string
	}{
		{
			name:   "input1",
			input:  schema.Operator(1),
			output: `"delete"`,
		},
	}

	for _, tt := range flagtests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := json.Marshal(tt.input)

			assert.NoError(t, err)
			assert.Equal(t, tt.output, string(res))
		})
	}
}

func TestOperation_OperatorZero(t *testing.T) {
	op := schema.Operator(0)
	input := schema.Operation{Operator: &op}

	assert.True(t, input.OperatorZero())
}

func TestOperation_OperatorDelete(t *testing.T) {
	op := schema.OperatorDelete
	input := schema.Operation{Operator: &op}

	assert.True(t, input.OperatorDelete())
}

func TestOperation_OperatorOrZero(t *testing.T) {
	input := schema.Operation{}

	assert.Equal(t, schema.Operator(0), input.OperatorOrZero())
}
