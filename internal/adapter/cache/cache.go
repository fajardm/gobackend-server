package cache

import (
	"context"
	"errors"
	"time"
)

const DefaultTTL time.Duration = 30 * time.Second

var ErrNil = errors.New("cache: value doesnt exists")

type Cache interface {
	Get(ctx context.Context, key string, output interface{}) error
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Del(ctx context.Context, keys ...string) error
	Clear(ctx context.Context) error
}
