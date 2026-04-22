package routes

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"

	"github.com/example/test-full-stack-developer/backend/docs"
	"github.com/example/test-full-stack-developer/backend/internal/adapter/handler"
	"github.com/example/test-full-stack-developer/backend/internal/adapter/repo"
	"github.com/example/test-full-stack-developer/backend/internal/config"
	"github.com/example/test-full-stack-developer/backend/internal/service"
)

func BindAPIRoutes(app *fiber.App, cfg config.Config, db *sql.DB) {
	imageRepo := repo.NewImageRepo(db)
	tagRepo := repo.NewTagRepo(db)
	imageTagRepo := repo.NewImageTagRepo(db)

	imageService := service.NewImageService(imageRepo, cfg.DefaultPageLimit, cfg.MaxPageLimit)
	tagService := service.NewTagService(tagRepo)
	imageTagService := service.NewImageTagService(imageTagRepo, imageRepo, tagRepo)

	imageHandler := handler.NewImageHandler(imageService)
	tagHandler := handler.NewTagHandler(tagService)
	imageTagHandler := handler.NewImageTagHandler(imageTagService)

	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})
	app.Get("/swagger", docs.SwaggerUI)
	app.Get("/swagger/", docs.SwaggerUI)
	app.Get("/swagger/openapi.yaml", docs.OpenAPI)

	api := app.Group(cfg.APIPrefix)
	BindImageRoutes(api.Group("/images"), imageHandler, imageTagHandler)
	BindTagRoutes(api.Group("/tags"), tagHandler)
}

func BindImageRoutes(router fiber.Router, imageHandler *handler.ImageHandler, imageTagHandler *handler.ImageTagHandler) {
	router.Post("/", imageHandler.Create)
	router.Get("/", imageHandler.List)
	router.Get("/:imageID", imageHandler.GetByID)
	router.Patch("/:imageID", imageHandler.Update)
	router.Delete("/:imageID", imageHandler.SoftDelete)
	router.Post("/:imageID/tags", imageTagHandler.Attach)
	router.Delete("/:imageID/tags/:tagID", imageTagHandler.Detach)
}

func BindTagRoutes(router fiber.Router, tagHandler *handler.TagHandler) {
	router.Post("/", tagHandler.Create)
	router.Get("/", tagHandler.List)
	router.Get("/:tagID", tagHandler.GetByID)
	router.Patch("/:tagID", tagHandler.Update)
	router.Delete("/:tagID", tagHandler.SoftDelete)
}
