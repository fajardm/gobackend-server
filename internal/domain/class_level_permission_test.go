package domain_test

import (
	"encoding/json"
	"testing"

	"github.com/fajardm/gobackend-server/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestClassLevelPermissions_UnmarshalJSON_Error(t *testing.T) {
	var flagtests = []struct {
		name   string
		input  string
		output domain.ClassLevelPermissions
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
			var value domain.ClassLevelPermissions
			err := json.Unmarshal([]byte(tt.input), &value)

			assert.Error(t, err)
		})
	}
}

func TestClassLevelPermissions_UnmarshalJSON_NoError(t *testing.T) {
	var flagtests = []struct {
		name   string
		input  string
		output domain.ClassLevelPermissions
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
			output: domain.ClassLevelPermissions{
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
				ProtectedFields: map[string]domain.FieldNames{
					"*":                                    domain.FieldNames{"name"},
					"role:admin":                           domain.FieldNames{"name"},
					"4590ccb6-43f1-45b7-9661-640b956d727c": domain.FieldNames{"name"},
				},
			},
		},
	}

	for _, tt := range flagtests {
		t.Run(tt.name, func(t *testing.T) {
			var value domain.ClassLevelPermissions
			err := json.Unmarshal([]byte(tt.input), &value)

			assert.NoError(t, err)
			assert.Equal(t, tt.output, value)
		})
	}
}
