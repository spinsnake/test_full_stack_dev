package docs

import (
	"embed"

	"github.com/gofiber/fiber/v2"
)

//go:embed openapi.yaml swagger-ui.html
var assets embed.FS

func SwaggerUI(c *fiber.Ctx) error {
	body, err := assets.ReadFile("swagger-ui.html")
	if err != nil {
		return fiber.ErrInternalServerError
	}

	c.Type("html", "utf-8")
	return c.Send(body)
}

func OpenAPI(c *fiber.Ctx) error {
	body, err := assets.ReadFile("openapi.yaml")
	if err != nil {
		return fiber.ErrInternalServerError
	}

	c.Type("yaml", "utf-8")
	return c.Send(body)
}
