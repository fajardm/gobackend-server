package handler

import (
	"github.com/fajardm/gobackend-server/internal/application"
	"github.com/gofiber/fiber/v2"
)

type Health struct {
	Service *application.Service `inject:"service"`
}

func (h Health) Ping(c *fiber.Ctx) error {
	return c.JSON(h.Service.Health.Ping())
}
