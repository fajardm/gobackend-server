package domain_test

import (
	"encoding/json"
	"testing"

	"github.com/fajardm/gobackend-server/internal/domain"
	"github.com/fajardm/gobackend-server/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestClass_UnmarshalJSON_Error(t *testing.T) {
	var flagtests = []struct {
		name   string
		input  string
		output errors.Code
	}{
		{
			name:   "input1",
			input:  `100`,
			output: errors.InvalidJSON,
		},
		{
			name: "input2",
			input: `
			{
				"fields":{
				   "_rperm":{
					  "type":"Array"
				   },
				   "_wperm":{
					  "type":"Array"
				   },
				   "objectId":{
					  "type":"UUID"
				   },
				   "createdAt":{
					  "type":"Date"
				   },
				   "updatedAt":{
					  "type":"Date"
				   }
				},
				"name":"Student",
				"classLevelPermissions":{
				   "get":{
					  "*":true
				   },
				   "find":{
					  "*":true
				   },
				   "count":{
					  "*":true
				   },
				   "create":{
					  "*":true
				   },
				   "delete":{
					  "*":true
				   },
				   "update":{
					  "*":true
				   },
				   "addField":{
					  "*":true
				   },
				   "protectedFields":{
					  "*":[
						 "xxx"
					  ]
				   }
				}
			}
			`,
			output: errors.InvalidJSON,
		},
		{
			name: "input3",
			input: `
			{
				"fields":{
				   "_rperm":{
					  "type":"Array"
				   },
				   "_wperm":{
					  "type":"Array"
				   },
				   "objectId":{
					  "type":"UUID"
				   },
				   "createdAt":{
					  "type":"Date"
				   },
				   "updatedAt":{
					  "type":"Date"
				   }
				},
				"name":"Student",
				"classLevelPermissions":{
				   "get":{
					  "*":true
				   },
				   "find":{
					  "*":true
				   },
				   "count":{
					  "*":true
				   },
				   "create":{
					  "*":true
				   },
				   "delete":{
					  "*":true
				   },
				   "update":{
					  "*":true
				   },
				   "addField":{
					  "*":true
				   }
				},
				"indexes":{
				   "xxx":{
					  "columns":[
						 "zzzz"
					  ]
				   }
				}
			}
			`,
			output: errors.InvalidJSON,
		},
	}

	for _, tt := range flagtests {
		t.Run(tt.name, func(t *testing.T) {
			var value domain.Class
			err := json.Unmarshal([]byte(tt.input), &value)

			actual := err.(errors.CustomError).Code()
			assert.Equal(t, tt.output, actual)
		})
	}
}

func TestClass_UnmarshalJSON_NoError(t *testing.T) {
	var flagtests = []struct {
		name   string
		input  string
		output domain.Class
	}{
		{
			name: "input1",
			input: `
			{
				"fields":{
				   "address": {
					  "type": "String"
				   }
				},
				"name":"_User"
			}
			`,
			output: domain.Class{
				Fields: domain.Fields{
					"address":       domain.Field{Type: domain.FieldTypeString, Required: false},
					"username":      domain.Field{Type: domain.FieldTypeString, Required: true},
					"password":      domain.Field{Type: domain.FieldTypeString, Required: true},
					"email":         domain.Field{Type: domain.FieldTypeString, Required: true},
					"emailVerified": domain.Field{Type: domain.FieldTypeBoolean, Required: true},
					"authData":      domain.Field{Type: domain.FieldTypeObject, Required: true},
					"_rperm":        domain.Field{Type: domain.FieldTypeArray, Required: false},
					"_wperm":        domain.Field{Type: domain.FieldTypeArray, Required: false},
					"objectId":      domain.Field{Type: domain.FieldTypeUUID, Required: false},
					"createdAt":     domain.Field{Type: domain.FieldTypeDate, Required: false},
					"updatedAt":     domain.Field{Type: domain.FieldTypeDate, Required: false},
				},
				Name: "_User",
				ClassLevelPermissions: domain.ClassLevelPermissions{
					Get:             map[string]bool{"*": true},
					Find:            map[string]bool{"*": true},
					Count:           map[string]bool{"*": true},
					Create:          map[string]bool{"*": true},
					Delete:          map[string]bool{"*": true},
					Update:          map[string]bool{"*": true},
					AddField:        map[string]bool{"*": true},
					ProtectedFields: map[string]domain.FieldNames{"*": make(domain.FieldNames, 0)},
				},
			},
		},
	}

	for _, tt := range flagtests {
		t.Run(tt.name, func(t *testing.T) {
			var value domain.Class
			err := json.Unmarshal([]byte(tt.input), &value)

			assert.NoError(t, err)
			assert.Equal(t, tt.output, value)
		})
	}
}
