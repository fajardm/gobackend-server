package schema

import (
	"encoding/json"
	"fmt"

	"github.com/fajardm/gobackend-server/pkg/errors"
)

type Field struct {
	Type         FieldType     `json:"type"`
	Required     bool          `json:"required"`
	DefaultValue *DefaultValue `json:"defaultValue,omitempty"`
	TargetClass  *ClassName    `json:"targetClass,omitempty"`
}

func (f Field) Validate() error {
	if _, ok := RelationFieldTypes[f.Type]; ok {
		if f.TargetClass == nil || f.TargetClass.String() == "" {
			return errors.New(errors.MissingRequiredField, fmt.Sprintf("field type %s needs a class name", f.Type))
		}
	}

	if f.DefaultValue != nil {
		if f.DefaultValue.Type() != f.Type {
			return errors.New(errors.IncorrectFieldType, "invalid field type and default value field type")
		}
	}

	if f.Required {
		if f.Type == FieldTypeRelation {
			return errors.New(errors.IncorrectFieldType, fmt.Sprintf(`The 'required' option is not applicable for %s`, FieldTypeRelation))
		}
	}
	return nil
}

func (f *Field) UnmarshalJSON(data []byte) error {
	type field Field
	if err := json.Unmarshal(data, (*field)(f)); err != nil {
		return errors.New(errors.InvalidJSON, err.Error())
	}
	return f.Validate()
}

type Fields map[FieldName]Field

func (f Fields) Delete(key FieldName) {
	delete(f, key)
}
