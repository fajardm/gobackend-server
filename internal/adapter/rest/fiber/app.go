package fiber

import (
	"github.com/fajardm/gobackend-server/config"
	"github.com/fajardm/gobackend-server/internal/adapter/rest/fiber/handler"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

type Handler struct {
	Health *handler.Health `inject:""`
}

type App struct {
	fiber   *fiber.App
	Config  *config.Config `inject:"config"`
	Handler *Handler       `inject:""`
}

func (a *App) Startup() error {
	a.fiber = fiber.New()
	go a.run()
	return nil
}

func (a App) Shutdown() error {
	return a.fiber.Shutdown()
}

func (a App) run() {
	// Init middleware
	a.fiber.Use(requestid.New())
	a.fiber.Use(recover.New())
	// Init router
	a.fiber.Get("/health/ping", a.Handler.Health.Ping)
	// a.Fiber.Post("/schema/create", timeout.New(a.Schema.Create, 5*time.Second))
	// Init server
	a.fiber.Listen(a.Config.App.Address)
}
