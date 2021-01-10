package domain

import (
	"encoding/json"
	"fmt"

	"github.com/fajardm/gobackend-server/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

type ClassLevelPermissions struct {
	Get             map[string]bool       `json:"get"`
	Find            map[string]bool       `json:"find"`
	Count           map[string]bool       `json:"count"`
	Create          map[string]bool       `json:"create"`
	Update          map[string]bool       `json:"update"`
	Delete          map[string]bool       `json:"delete"`
	AddField        map[string]bool       `json:"addField"`
	ProtectedFields map[string]FieldNames `json:"protectedFields"`
}

func (c ClassLevelPermissions) Validate() error {
	if err := c.ValidatePermissions(); err != nil {
		return err
	}
	if err := c.ValidateProtectedFields(); err != nil {
		return err
	}
	return nil
}

func (c *ClassLevelPermissions) UnmarshalJSON(data []byte) error {
	type classLevelPermissions ClassLevelPermissions
	if err := json.Unmarshal(data, (*classLevelPermissions)(c)); err != nil {
		return errors.New(errors.InvalidJSON, err.Error())
	}
	return c.Validate()
}

func (c ClassLevelPermissions) ValidatePermissions() error {
	operations := []map[string]bool{c.Get, c.Find, c.Count, c.Create, c.Update, c.Delete, c.AddField}
	for _, operation := range operations {
		for key, _ := range operation {
			if err := c.validatePermissionKey(key); err != nil {
				return err
			}
		}
	}
	return nil
}

func (c ClassLevelPermissions) ValidateProtectedFields() error {
	for key, protectedFields := range c.ProtectedFields {
		if err := c.validateProtectedFieldsKey(key); err != nil {
			return err
		}
		for _, field := range protectedFields {
			if _, ok := DefaultField[field]; ok {
				return errors.New(errors.InvalidJSON, fmt.Sprintf("default field %s can not be protected", key))
			}
		}
	}
	return nil
}

func (c ClassLevelPermissions) validateProtectedFieldsKey(key string) error {
	for _, regex := range ProtectedFieldsRegexes {
		if ok := regex.MatchString(key); ok {
			return nil
		}
	}
	if _, err := uuid.FromString(key); err == nil {
		return nil
	}
	return errors.New(errors.InvalidJSON, fmt.Sprintf("%s is not a valid key for class level permissions", key))
}

func (c ClassLevelPermissions) validatePermissionKey(key string) error {
	for _, regex := range ClassLevelPermissionsRegexes {
		if ok := regex.MatchString(key); ok {
			return nil
		}
	}
	if _, err := uuid.FromString(key); err == nil {
		return nil
	}
	return errors.New(errors.InvalidJSON, fmt.Sprintf("%s is not a valid key for class level permissions", key))
}
