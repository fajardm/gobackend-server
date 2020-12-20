package postgres

import (
	"context"

	"github.com/fajardm/gobackend-server/internal/domain/class"

	"github.com/fajardm/gobackend-server/internal/adapter/database"
	"github.com/fajardm/gobackend-server/internal/domain/schema"
)

type Class struct {
	DB        *database.Postgres `inject:"postgres"`
	Unmarshal Unmarshal          `inject:"unmarshal"`
}

func (r Class) Find(ctx context.Context, schema schema.Schema, query interface{}, queryOption class.QueryOption) error {
	// hasLimit, hasSkip := query.Limit > 0, query.Skip > 0
	return nil
}
