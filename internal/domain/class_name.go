package domain

import (
	"encoding/json"
	"fmt"

	"github.com/fajardm/gobackend-server/pkg/errors"
)

type ClassName string

func (c ClassName) String() string {
	return string(c)
}

func (c ClassName) Validate() error {
	if ok := JoinClassRegex.MatchString(c.String()); ok {
		return nil
	}
	if ok := ClassAndFieldRegex.MatchString(c.String()); ok {
		return nil
	}
	if _, ok := SystemClasses[c.String()]; ok {
		return nil
	}
	return errors.New(errors.InvalidClassName, fmt.Sprintf("%s classnames can only have alphanumeric characters and _, and must start with an alpha character", c))
}

func (c *ClassName) UnmarshalJSON(data []byte) error {
	type className ClassName
	if err := json.Unmarshal(data, (*className)(c)); err != nil {
		return errors.New(errors.InvalidJSON, err.Error())
	}
	return c.Validate()
}

func JoinClassName(fieldName FieldName, className ClassName) ClassName {
	return ClassName(fmt.Sprintf("_Join:%s:%s", fieldName.String(), className.String()))
}
