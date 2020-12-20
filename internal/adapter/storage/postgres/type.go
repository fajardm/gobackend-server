package postgres

import (
	"fmt"

	"github.com/fajardm/gobackend-server/internal/domain/schema"
)

type Unmarshal func(data []byte, v interface{}) error

func NewFieldType(fieldType schema.FieldType) (string, error) {
	switch fieldType {
	case schema.FieldTypeUUID:
		return "uuid", nil
	case schema.FieldTypeString:
		return "text", nil
	case schema.FieldTypeDate:
		return "timestamp with time zone", nil
	case schema.FieldTypeObject:
		return "jsonb", nil
	case schema.FieldTypeBoolean:
		return "boolean", nil
	case schema.FieldTypePointer:
		return "uuid", nil
	case schema.FieldTypeDecimal:
		return "double precision", nil
	case schema.FieldTypeArray:
		return "text[]", nil
	default:
		return "", fmt.Errorf("no type for %s yet", fieldType)
	}
}
