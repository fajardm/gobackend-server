package domain

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/fajardm/gobackend-server/pkg/errors"
)

type FieldType int

const (
	FieldTypeUUID FieldType = iota + 1
	FieldTypeBoolean
	FieldTypeString
	FieldTypeDecimal
	FieldTypeInteger
	FieldTypeDate
	FieldTypeArray
	FieldTypeObject
	FieldTypePointer
	FieldTypeRelation
)

var fieldTypeToString = map[FieldType]string{
	FieldTypeUUID:     "UUID",
	FieldTypeBoolean:  "Boolean",
	FieldTypeString:   "String",
	FieldTypeDecimal:  "Decimal",
	FieldTypeInteger:  "Integer",
	FieldTypeDate:     "Date",
	FieldTypeArray:    "Array",
	FieldTypeObject:   "Object",
	FieldTypePointer:  "Pointer",
	FieldTypeRelation: "Relation",
}

var stringToFieldType = map[string]FieldType{
	"UUID":     FieldTypeUUID,
	"Boolean":  FieldTypeBoolean,
	"String":   FieldTypeString,
	"Decimal":  FieldTypeDecimal,
	"Integer":  FieldTypeInteger,
	"Date":     FieldTypeDate,
	"Array":    FieldTypeArray,
	"Object":   FieldTypeObject,
	"Pointer":  FieldTypePointer,
	"Relation": FieldTypeRelation,
}

var reflectToFieldType = map[reflect.Kind]FieldType{
	reflect.Bool:    FieldTypeBoolean,
	reflect.String:  FieldTypeString,
	reflect.Float32: FieldTypeDecimal,
	reflect.Float64: FieldTypeDecimal,
	reflect.Int:     FieldTypeInteger,
	reflect.Int32:   FieldTypeInteger,
	reflect.Int64:   FieldTypeInteger,
}

var FieldTypes = map[FieldType]FieldType{
	FieldTypeUUID:     FieldTypeUUID,
	FieldTypeBoolean:  FieldTypeBoolean,
	FieldTypeString:   FieldTypeString,
	FieldTypeDecimal:  FieldTypeDecimal,
	FieldTypeInteger:  FieldTypeInteger,
	FieldTypeDate:     FieldTypeDate,
	FieldTypeArray:    FieldTypeArray,
	FieldTypeObject:   FieldTypeObject,
	FieldTypePointer:  FieldTypePointer,
	FieldTypeRelation: FieldTypeRelation,
}

var RelationFieldTypes = map[FieldType]FieldType{
	FieldTypePointer:  FieldTypePointer,
	FieldTypeRelation: FieldTypeRelation,
}

var NonRelationFieldTypes = map[FieldType]FieldType{
	FieldTypeUUID:    FieldTypeUUID,
	FieldTypeBoolean: FieldTypeBoolean,
	FieldTypeString:  FieldTypeString,
	FieldTypeDecimal: FieldTypeDecimal,
	FieldTypeInteger: FieldTypeInteger,
	FieldTypeDate:    FieldTypeDate,
	FieldTypeArray:   FieldTypeArray,
	FieldTypeObject:  FieldTypeObject,
}

func FieldTypeFromString(s string) (FieldType, error) {
	if o, ok := stringToFieldType[s]; ok {
		return o, nil
	}
	return 0, errors.New(errors.IncorrectFieldType, fmt.Sprintf("invalid operator %s", s))
}

func (f FieldType) ValidateValue(value interface{}) error {
	kind := reflect.TypeOf(value).Kind()
	if t, ok := reflectToFieldType[kind]; ok && t == f {
		return nil
	}
	return errors.New(errors.IncorrectFieldType, fmt.Sprintf("invalid field type of: %s", value))
}

func (f FieldType) String() string {
	return fieldTypeToString[f]
}

func (f FieldType) MarshalText() ([]byte, error) {
	s := f.String()
	if s == "" {
		return nil, errors.New(errors.IncorrectFieldType, fmt.Sprintf("invalid field type: %s", f))
	}
	return []byte(s), nil
}

func (f *FieldType) UnmarshalJSON(data []byte) error {
	var fieldType string
	if err := json.Unmarshal(data, &fieldType); err != nil {
		return errors.New(errors.InvalidJSON, err.Error())
	}
	ft, err := FieldTypeFromString(fieldType)
	if err != nil {
		return err
	}
	*f = ft
	return nil
}
