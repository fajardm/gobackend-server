package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fajardm/gobackend-server/config"
	"github.com/fajardm/gobackend-server/internal/adapter/cache"
	"github.com/fajardm/gobackend-server/internal/adapter/database"
	"github.com/fajardm/gobackend-server/internal/adapter/rest/fiber"
	"github.com/fajardm/gobackend-server/internal/adapter/storage/postgres"
	"github.com/fajardm/gobackend-server/internal/application"
	"github.com/fajardm/gobackend-server/internal/pkg/container"
)

func Serve(conf *config.Config) {
	log.Println("Serving...")

	container := container.New()

	container.RegisterService("config", conf)
	container.RegisterService("unmarshal", json.Unmarshal)
	container.RegisterService("fiber", new(fiber.App))
	container.RegisterService("cache", new(cache.Redis))
	container.RegisterService("service", new(application.Service))

	// Register Repository
	if conf.Database.Type == "postgres" {
		container.RegisterServices(map[string]interface{}{
			"postgres":         new(database.Postgres),
			"schemaQuerier":    postgres.GetSchemaQuerier(),
			"pingerRepository": new(postgres.Pinger),
			"schemaRepository": new(postgres.Schema),
		})
	}

	if err := container.Ready(); err != nil {
		log.Fatal("Failed to populate services", err)
	}

	quit := make(chan os.Signal, 2)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	delayShutdown := time.Duration(3) * time.Second

	log.Println(fmt.Sprintf("Signal termination received. Waiting %vs to shutdown.", delayShutdown.Seconds()))
	time.Sleep(delayShutdown)
	log.Println("Cleaning up resources...")

	container.Shutdown()

	log.Println("Bye")
}
