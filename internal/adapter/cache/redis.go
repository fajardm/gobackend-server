package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/fajardm/gobackend-server/config"
	"github.com/go-redis/redis/v8"
)

var _ Cache = (*Redis)(nil)

type Redis struct {
	Config *config.Config `inject:"config"`
	client *redis.Client
}

func (r *Redis) Startup() error {
	options := &redis.Options{
		Addr:     r.Config.Redis.Addresses[0],
		Password: r.Config.Redis.Password,
	}
	r.client = redis.NewClient(options)
	return nil
}

func (r Redis) Shutdown() error {
	return r.client.Close()
}

func (r Redis) Get(ctx context.Context, key string, output interface{}) error {
	values, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return ErrNil
	} else if err != nil {
		return err
	}
	return json.Unmarshal([]byte(values), &output)
}

func (r Redis) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}
	if ttl <= 0 {
		ttl = DefaultTTL
	}
	return r.client.Set(ctx, key, string(bytes), ttl).Err()
}

func (r Redis) Del(ctx context.Context, keys ...string) error {
	return r.client.Del(ctx, keys...).Err()
}

func (r Redis) Clear(ctx context.Context) error {
	return r.client.FlushAll(ctx).Err()
}
