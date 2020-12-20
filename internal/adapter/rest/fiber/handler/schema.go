package handler

import (
	"github.com/fajardm/gobackend-server/internal/application"
	"github.com/fajardm/gobackend-server/internal/domain/schema"
	"github.com/gofiber/fiber/v2"
)

type Schema struct {
	Service *application.Service `inject:"service"`
}

func (h Schema) Create(c *fiber.Ctx) error {
	payload := new(schema.Schema)
	if err := c.BodyParser(&payload); err != nil {
		return err
	}

	res, err := h.Service.Schema.Create(c.Context(), *payload)
	if err != nil {
		return err
	}

	return c.JSON(res)
}
