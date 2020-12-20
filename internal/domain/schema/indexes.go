package schema

import (
	"encoding/json"

	"github.com/fajardm/gobackend-server/pkg/errors"
)

type Index struct {
	Columns []FieldName `json:"columns"`
	Unique  bool        `json:"unique"`
}

func (i Index) Validate() error {
	if len(i.Columns) == 0 {
		return errors.New(errors.InvalidJSON, "index columns can not empty")
	}
	return nil
}

func (i *Index) UnmarshalJSON(data []byte) error {
	type index Index
	if err := json.Unmarshal(data, (*index)(i)); err != nil {
		return errors.New(errors.InvalidJSON, err.Error())
	}
	return i.Validate()
}

type Indexes map[string]Index

func (i Indexes) Delete(key string) {
	delete(i, key)
}
