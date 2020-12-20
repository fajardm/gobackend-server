package schema_test

import (
	"encoding/json"
	"testing"

	"github.com/fajardm/gobackend-server/internal/domain/schema"
	"github.com/stretchr/testify/assert"
)

func TestClassLevelPermissions_UnmarshalJSON_Error(t *testing.T) {
	var flagtests = []struct {
		name   string
		input  string
		output schema.ClassLevelPermissions
	}{
		{
			name: "invalidPermissionKey",
			input: `
			{
				"get":{
				   "xxxx":true
				},
				"protectedFields":{
				   "xxxx":[
					  "name"
				   ]
				}
			}
			`,
		},
		{
			name: "invalidProtectedKey",
			input: `
			{
				"get":{
				   "*":true
				},
				"protectedFields":{
				   "xxxx":[
					  "name"
				   ]
				}
			}
			`,
		},
		{
			name: "defaultFieldCanNotBeProtected",
			input: `
			{
				"get":{
				   "*":true
				},
				"protectedFields":{
				   "*":[
					  "objectId"
				   ]
				}
			}
			`,
		},
		{
			name: "invalidJSON",
			input: `
			{
				"get": "",
				"protectedFields":{
				   "*":[
					  "name"
				   ]
				}
			}
			`,
		},
	}

	for _, tt := range flagtests {
		t.Run(tt.name, func(t *testing.T) {
			var value schema.ClassLevelPermissions
			err := json.Unmarshal([]byte(tt.input), &value)

			assert.Error(t, err)
		})
	}
}

func TestClassLevelPermissions_UnmarshalJSON_NoError(t *testing.T) {
	var flagtests = []struct {
		name   string
		input  string
		output schema.ClassLevelPermissions
	}{
		{
			name: "all",
			input: `
			{
				"get":{
				   "*":true,
				   "role:admin":true,
				   "4590ccb6-43f1-45b7-9661-640b956d727c":true
				},
				"find":{
				   "*":true,
				   "role:admin":true,
				   "4590ccb6-43f1-45b7-9661-640b956d727c":true
				},
				"count":{
				   "*":true,
				   "role:admin":true,
				   "4590ccb6-43f1-45b7-9661-640b956d727c":true
				},
				"create":{
				   "*":true,
				   "role:admin":true,
				   "4590ccb6-43f1-45b7-9661-640b956d727c":true
				},
				"delete":{
				   "*":true,
				   "role:admin":true,
				   "4590ccb6-43f1-45b7-9661-640b956d727c":true
				},
				"update":{
				   "*":true,
				   "role:admin":true,
				   "4590ccb6-43f1-45b7-9661-640b956d727c":true
				},
				"addField":{
				   "*":true,
				   "role:admin":true,
				   "4590ccb6-43f1-45b7-9661-640b956d727c":true
				},
				"protectedFields":{
				   "*":[
					  "name"
				   ],
				   "role:admin":[
					  "name"
				   ],
				   "4590ccb6-43f1-45b7-9661-640b956d727c":[
					  "name"
				   ]
				}
			}
			`,
			output: schema.ClassLevelPermissions{
				Get: map[string]bool{
					"*":                                    true,
					"role:admin":                           true,
					"4590ccb6-43f1-45b7-9661-640b956d727c": true,
				},
				Find: map[string]bool{
					"*":                                    true,
					"role:admin":                           true,
					"4590ccb6-43f1-45b7-9661-640b956d727c": true,
				},
				Count: map[string]bool{
					"*":                                    true,
					"role:admin":                           true,
					"4590ccb6-43f1-45b7-9661-640b956d727c": true,
				},
				Create: map[string]bool{
					"*":                                    true,
					"role:admin":                           true,
					"4590ccb6-43f1-45b7-9661-640b956d727c": true,
				},
				Delete: map[string]bool{
					"*":                                    true,
					"role:admin":                           true,
					"4590ccb6-43f1-45b7-9661-640b956d727c": true,
				},
				Update: map[string]bool{
					"*":                                    true,
					"role:admin":                           true,
					"4590ccb6-43f1-45b7-9661-640b956d727c": true,
				},
				AddField: map[string]bool{
					"*":                                    true,
					"role:admin":                           true,
					"4590ccb6-43f1-45b7-9661-640b956d727c": true,
				},
				ProtectedFields: map[string]schema.FieldNames{
					"*":                                    schema.FieldNames{"name"},
					"role:admin":                           schema.FieldNames{"name"},
					"4590ccb6-43f1-45b7-9661-640b956d727c": schema.FieldNames{"name"},
				},
			},
		},
	}

	for _, tt := range flagtests {
		t.Run(tt.name, func(t *testing.T) {
			var value schema.ClassLevelPermissions
			err := json.Unmarshal([]byte(tt.input), &value)

			assert.NoError(t, err)
			assert.Equal(t, tt.output, value)
		})
	}
}
