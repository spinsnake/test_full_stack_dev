package httpx

import (
	"errors"
	"log"

	"github.com/gofiber/fiber/v2"

	"github.com/example/test-full-stack-developer/backend/internal/apperrors"
)

func Data(c *fiber.Ctx, status int, payload any) error {
	return c.Status(status).JSON(fiber.Map{"data": payload})
}

func Message(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(fiber.Map{"message": message})
}

func Error(c *fiber.Ctx, err error) error {
	status := fiber.StatusInternalServerError
	message := "internal server error"

	var fiberError *fiber.Error
	if errors.As(err, &fiberError) {
		status = fiberError.Code
		message = fiberError.Message
		return c.Status(status).JSON(fiber.Map{"error": message})
	}

	switch {
	case errors.Is(err, apperrors.ErrInvalidInput):
		status = fiber.StatusBadRequest
		message = err.Error()
	case errors.Is(err, apperrors.ErrNotFound):
		status = fiber.StatusNotFound
		message = err.Error()
	case errors.Is(err, apperrors.ErrConflict):
		status = fiber.StatusConflict
		message = err.Error()
	}

	if status >= fiber.StatusInternalServerError {
		log.Printf("http error: %v", err)
	}

	return c.Status(status).JSON(fiber.Map{"error": message})
}
