package domain

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/fajardm/gobackend-server/pkg/errors"
)

type ToUpdateClass struct {
	AddFields     Fields
	DeleteFields  Fields
	AddIndexes    Indexes
	DeleteIndexes Indexes
}

type Class struct {
	Name                  ClassName             `json:"name"`
	Fields                Fields                `json:"fields"`
	Indexes               Indexes               `json:"indexes"`
	ClassLevelPermissions ClassLevelPermissions `json:"classLevelPermissions"`
	Objects               Objects               `json:"objects"`
	ToUpdate              *ToUpdateClass        `json:"toUpdate"`
}

func (c Class) Validate() error {
	if err := c.ValidateExistenceProtectedFields(); err != nil {
		return err
	}
	if err := c.ValidateExistenceIndexField(); err != nil {
		return err
	}
	if err := c.ValidateObjects(); err != nil {
		return err
	}
	return nil
}

func (c Class) ValidateExistenceProtectedFields() error {
	for key, protectedFields := range c.ClassLevelPermissions.ProtectedFields {
		for _, field := range protectedFields {
			if _, ok := c.Fields[field]; !ok {
				return errors.New(errors.InvalidJSON, fmt.Sprintf("field '%s' in protectedFields:%s does not exist", field, key))
			}
		}
	}
	return nil
}

func (c Class) ValidateExistenceIndexField() error {
	for _, index := range c.Indexes {
		for _, column := range index.Columns {
			if _, ok := c.Fields[column]; !ok {
				return errors.New(errors.InvalidJSON, fmt.Sprintf(`field %s does not exists, cannot add index`, column))
			}
		}
	}
	return nil
}

func (c Class) AddDefaultFields() {
	for name, field := range DefaultField {
		c.Fields[name] = field
	}

	if defaultFields, ok := DefaultFields[c.Name.String()]; ok {
		for name, field := range defaultFields {
			c.Fields[name] = field
		}
	}
}

func (c *Class) AddDefaultClassLevelPermissions() {
	c.ClassLevelPermissions = ClassLevelPermissions{
		Get:             map[string]bool{"*": true},
		Find:            map[string]bool{"*": true},
		Count:           map[string]bool{"*": true},
		Create:          map[string]bool{"*": true},
		Delete:          map[string]bool{"*": true},
		Update:          map[string]bool{"*": true},
		AddField:        map[string]bool{"*": true},
		ProtectedFields: map[string]FieldNames{"*": make(FieldNames, 0)},
	}
}

func (c Class) ValidateObjects() error {
	for _, object := range c.Objects {
		for key, val := range object {
			field, ok := c.Fields[key]
			if !ok {
				return errors.New(errors.InvalidJSON, fmt.Sprintf("field '%s' does not exist", key))
			}
			if err := field.Type.ValidateValue(val); err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *Class) Update(newClass Class) {
	var wg sync.WaitGroup
	c.ToUpdate = new(ToUpdateClass)

	go func() {
		wg.Add(1)
		defer wg.Done()
		for key, field := range newClass.Fields {
			if _, ok := c.Fields[key]; !ok {
				c.Fields[key] = field
				c.ToUpdate.AddFields[key] = field
			}
		}
	}()

	go func() {
		wg.Add(1)
		defer wg.Done()
		for key, field := range c.Fields {
			if _, ok := newClass.Fields[key]; !ok {
				c.Fields.Delete(key)
				c.ToUpdate.DeleteFields[key] = field
			}
		}
	}()

	go func() {
		wg.Add(1)
		defer wg.Done()
		for key, index := range newClass.Indexes {
			if _, ok := c.Indexes[key]; !ok {
				c.Indexes[key] = index
				c.ToUpdate.AddIndexes[key] = index
			}
		}
	}()

	go func() {
		wg.Add(1)
		defer wg.Done()
		for key, index := range c.Indexes {
			if _, ok := newClass.Indexes[key]; !ok {
				c.Indexes.Delete(key)
				c.ToUpdate.DeleteIndexes[key] = index
			}
		}
	}()

	wg.Wait()
}

func (c *Class) UnmarshalJSON(data []byte) error {
	type schema Class
	if err := json.Unmarshal(data, (*schema)(c)); err != nil {
		return errors.New(errors.InvalidJSON, err.Error())
	}
	if err := c.Validate(); err != nil {
		return err
	}
	c.AddDefaultFields()
	c.AddDefaultClassLevelPermissions()
	return nil
}

type Classes map[ClassName]Class
