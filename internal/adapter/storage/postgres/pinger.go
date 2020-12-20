package postgres

import "github.com/fajardm/gobackend-server/internal/adapter/database"

type Pinger struct {
	DB *database.Postgres `inject:"postgres"`
}

func (r Pinger) Ping() error {
	return r.DB.Ping()
}
