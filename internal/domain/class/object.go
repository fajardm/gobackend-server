package class

import (
	"encoding/json"
	"fmt"

	"github.com/fajardm/gobackend-server/internal/domain/schema"
	"github.com/fajardm/gobackend-server/pkg/errors"
)

type ObjectValue map[schema.FieldName]interface{}

type Object struct {
	value  ObjectValue
	schema schema.Schema
}

func NewObject(schema schema.Schema) *Object {
	return &Object{schema: schema}
}

// NewObjectWithValue creates a new instance just for testing
func NewObjectWithValue(schema schema.Schema, value ObjectValue) *Object {
	return &Object{schema: schema, value: value}
}

func (o Object) Value() ObjectValue {
	return o.value
}

func (o Object) Validate() error {
	for key, val := range o.value {
		field, ok := o.schema.Fields[key]
		if !ok {
			return errors.New(errors.InvalidJSON, fmt.Sprintf("field '%s' does not exist", key))
		}
		if err := field.Type.ValidateValue(val); err != nil {
			return err
		}
	}
	return nil
}

func (o Object) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.value)
}

func (o *Object) UnmarshalJSON(data []byte) error {
	var objectValue ObjectValue
	if err := json.Unmarshal(data, &objectValue); err != nil {
		return errors.New(errors.InvalidJSON, err.Error())
	}
	o.value = objectValue
	return o.Validate()
}
