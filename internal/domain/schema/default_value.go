package schema

import (
	"encoding/json"
	"time"

	"github.com/fajardm/gobackend-server/pkg/errors"
)

type DefaultValuer interface {
	FieldType() FieldType
}

type DefaultValuePointer struct {
	Type        FieldType `json:"type"`
	TargetClass ClassName `json:"className"`
	Value       string    `json:"value"`
}

func (d DefaultValuePointer) FieldType() FieldType {
	return d.Type
}

type DefaultValueDate struct {
	Type  FieldType `json:"type"`
	Value time.Time `json:"value"`
}

func (d DefaultValueDate) FieldType() FieldType {
	return d.Type
}

type DefaultValueString struct {
	Type  FieldType `json:"type"`
	Value string    `json:"value"`
}

func (d DefaultValueString) FieldType() FieldType {
	return d.Type
}

type DefaultValueInteger struct {
	Type  FieldType `json:"type"`
	Value int64     `json:"value"`
}

func (d DefaultValueInteger) FieldType() FieldType {
	return d.Type
}

type DefaultValueDecimal struct {
	Type  FieldType `json:"type"`
	Value float64   `json:"value"`
}

func (d DefaultValueDecimal) FieldType() FieldType {
	return d.Type
}

type DefaultValueBool struct {
	Type  FieldType `json:"type"`
	Value bool      `json:"value"`
}

func (d DefaultValueBool) FieldType() FieldType {
	return d.Type
}

type DefaultValueObject struct {
	Type  FieldType              `json:"type"`
	Value map[string]interface{} `json:"value"`
}

func (d DefaultValueObject) FieldType() FieldType {
	return d.Type
}

type DefaultValueArray struct {
	Type  FieldType     `json:"type"`
	Value []interface{} `json:"value"`
}

func (d DefaultValueArray) FieldType() FieldType {
	return d.Type
}

type DefaultValue struct {
	value interface{}
}

func NewDefaultValue(v interface{}) *DefaultValue {
	return &DefaultValue{v}
}

func (d DefaultValue) Type() FieldType {
	if v, ok := d.value.(DefaultValuer); ok {
		return v.FieldType()
	}
	return 0
}

func (d DefaultValue) Value() interface{} {
	return d.value
}

func (d *DefaultValue) UnmarshalJSON(data []byte) error {
	var Type struct {
		Type FieldType `json:"type"`
	}

	if err := json.Unmarshal(data, &Type); err != nil {
		return errors.New(errors.InvalidJSON, err.Error())
	}

	switch Type.Type {
	case FieldTypeBoolean:
		var value DefaultValueBool
		if err := json.Unmarshal(data, &value); err != nil {
			return errors.New(errors.InvalidJSON, err.Error())
		}
		d.value = value
	case FieldTypeString:
		var value DefaultValueString
		if err := json.Unmarshal(data, &value); err != nil {
			return errors.New(errors.InvalidJSON, err.Error())
		}
		d.value = value
	case FieldTypeDecimal:
		var value DefaultValueDecimal
		if err := json.Unmarshal(data, &value); err != nil {
			return errors.New(errors.InvalidJSON, err.Error())
		}
		d.value = value
	case FieldTypeInteger:
		var value DefaultValueInteger
		if err := json.Unmarshal(data, &value); err != nil {
			return errors.New(errors.InvalidJSON, err.Error())
		}
		d.value = value
	case FieldTypeDate:
		var value DefaultValueDate
		if err := json.Unmarshal(data, &value); err != nil {
			return errors.New(errors.InvalidJSON, err.Error())
		}
		d.value = value
	case FieldTypeArray:
		var value DefaultValueArray
		if err := json.Unmarshal(data, &value); err != nil {
			return errors.New(errors.InvalidJSON, err.Error())
		}
		d.value = value
	case FieldTypeObject:
		var value DefaultValueObject
		if err := json.Unmarshal(data, &value); err != nil {
			return errors.New(errors.InvalidJSON, err.Error())
		}
		d.value = value
	case FieldTypePointer:
		var value DefaultValuePointer
		if err := json.Unmarshal(data, &value); err != nil {
			return errors.New(errors.InvalidJSON, err.Error())
		}
		d.value = value
	}
	return nil
}
