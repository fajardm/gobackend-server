package service

import (
	"github.com/fajardm/gobackend-server/internal/application/model"
)

type Pinger interface {
	Ping() error
}

type Health struct {
	Pinger Pinger `inject:"pingerRepository"`
}

func (s Health) Ping() *model.Health {
	return &model.Health{
		Database: s.Pinger.Ping() == nil,
	}
}
