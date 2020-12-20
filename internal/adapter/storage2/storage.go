package storage2

import (
	"context"

	"github.com/fajardm/gobackend-server/internal/domain/schema"
)

type Storage interface {
	ClassExists(ctx context.Context, className schema.ClassName) (bool, error)
}
