package postgres

import (
	"fmt"

	"github.com/fajardm/gobackend-server/internal/domain"
)

type Unmarshal func(data []byte, v interface{}) error

func NewFieldType(fieldType domain.FieldType) (string, error) {
	switch fieldType {
	case domain.FieldTypeUUID:
		return "uuid", nil
	case domain.FieldTypeString:
		return "text", nil
	case domain.FieldTypeDate:
		return "timestamp with time zone", nil
	case domain.FieldTypeObject:
		return "jsonb", nil
	case domain.FieldTypeBoolean:
		return "boolean", nil
	case domain.FieldTypePointer:
		return "uuid", nil
	case domain.FieldTypeDecimal:
		return "double precision", nil
	case domain.FieldTypeArray:
		return "text[]", nil
	default:
		return "", fmt.Errorf("no type for %s yet", fieldType)
	}
}
