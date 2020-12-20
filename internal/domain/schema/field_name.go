package schema

import (
	"encoding/json"
	"fmt"

	"github.com/fajardm/gobackend-server/pkg/errors"
)

type FieldName string

func (f FieldName) String() string {
	return string(f)
}

func (f FieldName) Validate() error {
	if ok := ClassAndFieldRegex.MatchString(f.String()); !ok {
		return errors.New(errors.InvalidFieldName, fmt.Sprintf("invalid field name: %s", f))
	}
	return nil
}

func (f *FieldName) UnmarshalJSON(data []byte) error {
	type fieldName FieldName
	if err := json.Unmarshal(data, (*fieldName)(f)); err != nil {
		return errors.New(errors.InvalidJSON, err.Error())
	}
	return f.Validate()
}

const (
	FieldObjectID        FieldName = "objectId"
	FieldCreatedAt                 = "createdAt"
	FieldUpdatedAt                 = "updatedAt"
	FieldReadPermission            = "_rperm"
	FieldWritePermission           = "_wperm"
)

type FieldNames []FieldName
