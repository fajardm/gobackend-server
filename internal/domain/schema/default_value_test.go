package schema_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/fajardm/gobackend-server/internal/domain/schema"
	"github.com/stretchr/testify/assert"
)

func TestDefaultValue_UnmarshalJSON_Error(t *testing.T) {
	type Data struct {
		DefaultValue schema.DefaultValue `json:"defaultValue"`
	}

	var flagtests = []struct {
		name     string
		json     string
		expected Data
	}{
		{
			name: "Unknown",
			json: `{"defaultValue": {"type": "Unknown", "value": "xxxx"}}`,
		},
		{
			name: "Boolean",
			json: `{"defaultValue": {"type": "Boolean", "value": "true"}}`,
		},
		{
			name: "String",
			json: `{"defaultValue": {"type": "String", "value": 100}}`,
		},
		{
			name: "Decimal",
			json: `{"defaultValue": {"type": "Decimal", "value": "100.87212"}}`,
		},
		{
			name: "Integer",
			json: `{"defaultValue": {"type": "Integer", "value": "100"}}`,
		},
		{
			name: "Date",
			json: `{"defaultValue": {"type": "Date", "value": 100}}`,
		},
		{
			name: "Array",
			json: `{"defaultValue": {"type": "Array", "value": 100}}`,
		},
		{
			name: "Object",
			json: `{"defaultValue": {"type": "Object", "value": [{"name": "xxxx", "address": {"city": "xxxx"}}]}}`,
		},
		{
			name: "Pointer",
			json: `{"defaultValue": {"type": "Pointer", "className": "User", "value": 100}}`,
		},
	}

	for _, tt := range flagtests {
		t.Run(tt.name, func(t *testing.T) {
			var data Data
			err := json.Unmarshal([]byte(tt.json), &data)

			assert.Error(t, err)
		})
	}
}

func TestDefaultValue_UnmarshalJSON_NoError(t *testing.T) {
	now := time.Now()
	nowString := now.Format(time.RFC3339)
	nowTime, _ := time.Parse(time.RFC3339, nowString)

	type Data struct {
		DefaultValue schema.DefaultValue `json:"defaultValue"`
	}

	var flagtests = []struct {
		name     string
		json     string
		expected Data
	}{
		{
			name:     "Boolean",
			json:     `{"defaultValue": {"type": "Boolean", "value": true}}`,
			expected: Data{DefaultValue: *schema.NewDefaultValue(schema.DefaultValueBool{Type: schema.FieldTypeBoolean, Value: true})},
		},
		{
			name:     "String",
			json:     `{"defaultValue": {"type": "String", "value": "xxxx"}}`,
			expected: Data{DefaultValue: *schema.NewDefaultValue(schema.DefaultValueString{Type: schema.FieldTypeString, Value: "xxxx"})},
		},
		{
			name:     "Decimal",
			json:     `{"defaultValue": {"type": "Decimal", "value": 100.87212}}`,
			expected: Data{DefaultValue: *schema.NewDefaultValue(schema.DefaultValueDecimal{Type: schema.FieldTypeDecimal, Value: 100.87212})},
		},
		{
			name:     "Integer",
			json:     `{"defaultValue": {"type": "Integer", "value": 100}}`,
			expected: Data{DefaultValue: *schema.NewDefaultValue(schema.DefaultValueInteger{Type: schema.FieldTypeInteger, Value: 100})},
		},
		{
			name:     "Date",
			json:     fmt.Sprintf(`{"defaultValue": {"type": "Date", "value": "%s"}}`, nowString),
			expected: Data{DefaultValue: *schema.NewDefaultValue(schema.DefaultValueDate{Type: schema.FieldTypeDate, Value: nowTime})},
		},
		{
			name:     "Array",
			json:     `{"defaultValue": {"type": "Array", "value": [1,2,3,4]}}`,
			expected: Data{DefaultValue: *schema.NewDefaultValue(schema.DefaultValueArray{Type: schema.FieldTypeArray, Value: []interface{}{float64(1), float64(2), float64(3), float64(4)}})},
		},
		{
			name:     "Object",
			json:     `{"defaultValue": {"type": "Object", "value": {"name": "xxxx", "address": {"city": "xxxx"}}}}`,
			expected: Data{DefaultValue: *schema.NewDefaultValue(schema.DefaultValueObject{Type: schema.FieldTypeObject, Value: map[string]interface{}{"name": "xxxx", "address": map[string]interface{}{"city": "xxxx"}}})},
		},
		{
			name:     "Pointer",
			json:     `{"defaultValue": {"type": "Pointer", "className": "User", "value": "xxxx"}}`,
			expected: Data{DefaultValue: *schema.NewDefaultValue(schema.DefaultValuePointer{Type: schema.FieldTypePointer, TargetClass: "User", Value: "xxxx"})},
		},
	}

	for _, tt := range flagtests {
		t.Run(tt.name, func(t *testing.T) {
			var data Data
			err := json.Unmarshal([]byte(tt.json), &data)

			assert.NoError(t, err)
			assert.Equal(t, tt.expected, data)
			assert.Equal(t, tt.name, data.DefaultValue.Type().String())
		})
	}
}

func TestDefaultValue_Type_EmptyString(t *testing.T) {
	dv := schema.NewDefaultValue(100)

	assert.Equal(t, "", dv.Type().String())
}

func TestDefaultValue_Type_Value(t *testing.T) {
	dv := schema.NewDefaultValue("xxxx")

	assert.IsType(t, "string", dv.Value())
}
