package application

import "github.com/fajardm/gobackend-server/internal/application/service"

type Service struct {
	Health *service.Health `inject:""`
	Schema *service.Schema `inject:""`
}
