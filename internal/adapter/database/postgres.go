package database

import (
	"fmt"

	"github.com/fajardm/gobackend-server/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Postgres struct {
	Config   *config.Config `inject:"config"`
	*sqlx.DB `inject:""`
}

func (d *Postgres) Startup() error {
	conf := d.Config.Database.SQL
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", conf.Host, conf.Port, conf.Username, conf.Password, conf.Database)
	d.DB = sqlx.MustConnect("postgres", dsn).Unsafe()
	return nil
}

func (d Postgres) Shutdown() error {
	return d.Close()
}
