package schema_test

import (
	"encoding/json"
	"testing"

	"github.com/fajardm/gobackend-server/internal/domain/schema"
	"github.com/fajardm/gobackend-server/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestSchema_UnmarshalJSON_Error(t *testing.T) {
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
				"className":"Student",
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
				"className":"Student",
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
			var value schema.Schema
			err := json.Unmarshal([]byte(tt.input), &value)

			actual := err.(errors.CustomError).Code()
			assert.Equal(t, tt.output, actual)
		})
	}
}

func TestSchema_UnmarshalJSON_NoError(t *testing.T) {
	var flagtests = []struct {
		name   string
		input  string
		output schema.Schema
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
				"className":"CustomUser",
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
						 
					  ]
				   }
				}
			}
			`,
			output: schema.Schema{
				Fields: schema.Fields{
					"address":   schema.Field{Type: schema.FieldTypeString, Required: false},
					"_rperm":    schema.Field{Type: schema.FieldTypeArray, Required: false},
					"_wperm":    schema.Field{Type: schema.FieldTypeArray, Required: false},
					"objectId":  schema.Field{Type: schema.FieldTypeUUID, Required: false},
					"createdAt": schema.Field{Type: schema.FieldTypeDate, Required: false},
					"updatedAt": schema.Field{Type: schema.FieldTypeDate, Required: false},
				},
				ClassName: "CustomUser",
				ClassLevelPermissions: schema.ClassLevelPermissions{
					Get:             map[string]bool{"*": true},
					Find:            map[string]bool{"*": true},
					Count:           map[string]bool{"*": true},
					Create:          map[string]bool{"*": true},
					Delete:          map[string]bool{"*": true},
					Update:          map[string]bool{"*": true},
					AddField:        map[string]bool{"*": true},
					ProtectedFields: map[string]schema.FieldNames{"*": make(schema.FieldNames, 0)},
				},
			},
		},
	}

	for _, tt := range flagtests {
		t.Run(tt.name, func(t *testing.T) {
			var value schema.Schema
			err := json.Unmarshal([]byte(tt.input), &value)

			assert.NoError(t, err)
			assert.Equal(t, tt.output, value)
		})
	}
}
