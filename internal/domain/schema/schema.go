package schema

import (
	"encoding/json"
	"fmt"

	"github.com/fajardm/gobackend-server/pkg/errors"
)

type Schema struct {
	ClassName             ClassName             `json:"className"`
	Fields                Fields                `json:"fields"`
	Indexes               Indexes               `json:"indexes,omitempty"`
	ClassLevelPermissions ClassLevelPermissions `json:"classLevelPermissions"`
}

func (s Schema) Validate() error {
	if err := s.validateExistenceProtectedFields(); err != nil {
		return err
	}
	if err := s.validateExistenceIndexField(); err != nil {
		return err
	}
	return nil
}

func (s Schema) validateExistenceProtectedFields() error {
	for key, protectedFields := range s.ClassLevelPermissions.ProtectedFields {
		for _, field := range protectedFields {
			if _, ok := s.Fields[field]; !ok {
				return errors.New(errors.InvalidJSON, fmt.Sprintf("field '%s' in protectedFields:%s does not exist", field, key))
			}
		}
	}
	return nil
}

func (s Schema) validateExistenceIndexField() error {
	for _, index := range s.Indexes {
		for _, column := range index.Columns {
			if _, ok := s.Fields[column]; !ok {
				return errors.New(errors.InvalidJSON, fmt.Sprintf(`field %s does not exists, cannot add index`, column))
			}
		}
	}
	return nil
}

func (s Schema) addDefaultColumns() {
	for name, field := range DefaultColumn {
		s.Fields[name] = field
	}

	if defaultColumns, ok := DefaultColumns[s.ClassName.String()]; ok {
		for name, field := range defaultColumns {
			s.Fields[name] = field
		}
	}
}

func (s *Schema) UnmarshalJSON(data []byte) error {
	type schema Schema
	if err := json.Unmarshal(data, (*schema)(s)); err != nil {
		return errors.New(errors.InvalidJSON, err.Error())
	}
	if err := s.Validate(); err != nil {
		return err
	}
	s.addDefaultColumns()
	return nil
}

type Schemas map[ClassName]Schema
