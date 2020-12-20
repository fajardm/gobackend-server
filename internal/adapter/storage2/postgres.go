package storage2

import (
	"context"
	"fmt"

	"github.com/fajardm/gobackend-server/config"
	"github.com/fajardm/gobackend-server/internal/domain/schema"
	"github.com/jmoiron/sqlx"
)

type Unmarshal func(data []byte, v interface{}) error

type Postgres struct {
	Config    *config.Config `inject:"config"`
	Unmarshal Unmarshal      `inject:"unmarshal"`
	querier   PostgresQuerier
	client    *sqlx.DB
}

func (s *Postgres) Startup() error {
	conf := s.Config.Database.SQL
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", conf.Host, conf.Port, conf.Username, conf.Password, conf.Database)
	s.client = sqlx.MustConnect("postgres", dsn).Unsafe()
	return nil
}

func (s Postgres) Shutdown() error {
	return s.client.Close()
}

func (s Postgres) ClassExists(ctx context.Context, className schema.ClassName) (bool, error) {
	query, args, err := s.querier.ClassExists(className)
	if err != nil {
		return false, err
	}

	var res bool
	if err := s.client.QueryRowxContext(ctx, query, args...).Scan(&res); err != nil {
		return false, err
	}

	return res, nil
}
