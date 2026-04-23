package main

import (
	"log"
	"os"

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

	if len(os.Args) > 1 && os.Args[1] == "migrate" {
		migrateAction := "up"
		if len(os.Args) > 2 {
			migrateAction = os.Args[2]
		}
		if migrateAction != "up" {
			log.Fatalf("unsupported migrate action: %s", migrateAction)
		}

		if err := infra.RunMigrations(db, "", cfg.MockData); err != nil {
			log.Fatalf("run migrations: %v", err)
		}

		log.Print("migrations applied")
		return
	}

	if cfg.AutoMigrate {
		if err := infra.RunMigrations(db, "", cfg.MockData); err != nil {
			log.Fatalf("auto migrate: %v", err)
		}
		log.Print("auto migrations applied")
	}

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
