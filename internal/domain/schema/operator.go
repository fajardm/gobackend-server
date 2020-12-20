package schema

import (
	"encoding/json"
	"fmt"

	"github.com/fajardm/gobackend-server/pkg/errors"
)

type Operator int

const (
	OperatorDelete Operator = iota + 1
)

func OperatorFromString(s string) (Operator, error) {
	switch s {
	case "delete":
		return OperatorDelete, nil
	default:
		return 0, errors.New(errors.IncorrectOperation, fmt.Sprintf("invalid operator %s", s))
	}
}

func (o Operator) String() string {
	return [...]string{"", "delete"}[o]
}

func (o Operator) MarshalText() ([]byte, error) {
	s := o.String()
	if s == "" {
		return nil, errors.New(errors.IncorrectOperation, "invalid operator")
	}
	return []byte(s), nil
}

func (o *Operator) UnmarshalJSON(data []byte) error {
	var operator string
	if err := json.Unmarshal(data, &operator); err != nil {
		return errors.New(errors.InvalidJSON, err.Error())
	}
	op, err := OperatorFromString(operator)
	if err != nil {
		return err
	}
	*o = op
	return nil
}

type Operation struct {
	Operator *Operator `json:"__op,omitempty"`
}

func (o Operation) OperatorOrZero() Operator {
	if o.Operator == nil {
		return Operator(0)
	}
	return *o.Operator
}

func (o Operation) OperatorDelete() bool {
	return o.OperatorOrZero() == OperatorDelete
}

func (o Operation) OperatorZero() bool {
	return o.OperatorOrZero() == Operator(0)
}
