package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/example/test-full-stack-developer/backend/internal/config"
	"github.com/example/test-full-stack-developer/backend/internal/infra"
	"github.com/example/test-full-stack-developer/backend/internal/routes"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	db, err := infra.OpenMySQL(cfg)
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	defer db.Close()

	app := fiber.New(fiber.Config{
		AppName:      cfg.AppName,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CORSAllowOrigins,
		AllowMethods:     cfg.CORSAllowMethods,
		AllowHeaders:     cfg.CORSAllowHeaders,
		AllowCredentials: cfg.CORSAllowCredentials,
	}))

	routes.BindAPIRoutes(app, cfg, db)

	addr := cfg.Address()
	log.Printf("listening on %s", addr)

	if err := app.Listen(addr); err != nil {
		log.Fatalf("listen: %v", err)
	}
}
