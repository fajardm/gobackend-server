package service

import (
	"context"
	"fmt"

	"github.com/fajardm/gobackend-server/internal/domain/schema"
	"github.com/fajardm/gobackend-server/pkg/errors"
)

var _ schema.Service = (*Schema)(nil)

type Schema struct {
	Repository schema.Repository `inject:"schemaRepository"`
}

func (s Schema) FindByClassName(ctx context.Context, className schema.ClassName) (*schema.Schema, error) {
	return s.Repository.FindByClassName(ctx, className)
}

func (s Schema) Create(ctx context.Context, data schema.Schema) (*schema.Schema, error) {
	exists, err := s.Repository.Exists(ctx, data.ClassName)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New(errors.DuplicateData, fmt.Sprintf("class %s already exists", data.ClassName))
	}

	if err := s.Repository.Create(ctx, data); err != nil {
		return nil, err
	}

	return s.Repository.FindByClassName(ctx, data.ClassName)
}

func (s Schema) Update(ctx context.Context, data schema.Schema) (*schema.Schema, error) {
	exists, err := s.Repository.Exists(ctx, data.ClassName)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New(errors.DataNotFound, fmt.Sprintf("class %s not exists", data.ClassName))
	}

	if err := s.Repository.Update(ctx, data); err != nil {
		return nil, err
	}

	return s.Repository.FindByClassName(ctx, data.ClassName)
}

func (s Schema) Delete(ctx context.Context, className schema.ClassName) error {
	return s.Repository.Delete(ctx, className)
}
