package handler

import (
	"github.com/gofiber/fiber/v2"

	"github.com/example/test-full-stack-developer/backend/internal/entities"
	"github.com/example/test-full-stack-developer/backend/internal/httpx"
	"github.com/example/test-full-stack-developer/backend/internal/port"
)

type ImageHandler struct {
	service port.ImageService
}

func NewImageHandler(service port.ImageService) *ImageHandler {
	return &ImageHandler{service: service}
}

func (h *ImageHandler) Create(c *fiber.Ctx) error {
	var input entities.CreateImageInput
	if err := c.BodyParser(&input); err != nil {
		return httpx.Error(c, err)
	}

	result, err := h.service.Create(c.Context(), input)
	if err != nil {
		return httpx.Error(c, err)
	}

	return httpx.Data(c, fiber.StatusCreated, result)
}

func (h *ImageHandler) List(c *fiber.Ctx) error {
	cursor, err := parseOptionalUint64(c.Query("cursor"))
	if err != nil {
		return httpx.Error(c, err)
	}

	limit, err := parseOptionalInt(c.Query("limit"))
	if err != nil {
		return httpx.Error(c, err)
	}

	result, err := h.service.List(c.Context(), cursor, limit, c.Query("tag"))
	if err != nil {
		return httpx.Error(c, err)
	}

	return httpx.Data(c, fiber.StatusOK, result)
}

func (h *ImageHandler) GetByID(c *fiber.Ctx) error {
	imageID, err := parseRequiredUint64(c.Params("imageID"), "imageID")
	if err != nil {
		return httpx.Error(c, err)
	}

	result, err := h.service.GetByID(c.Context(), imageID)
	if err != nil {
		return httpx.Error(c, err)
	}

	return httpx.Data(c, fiber.StatusOK, result)
}

func (h *ImageHandler) Update(c *fiber.Ctx) error {
	imageID, err := parseRequiredUint64(c.Params("imageID"), "imageID")
	if err != nil {
		return httpx.Error(c, err)
	}

	var input entities.UpdateImageInput
	if err := c.BodyParser(&input); err != nil {
		return httpx.Error(c, err)
	}

	result, err := h.service.Update(c.Context(), imageID, input)
	if err != nil {
		return httpx.Error(c, err)
	}

	return httpx.Data(c, fiber.StatusOK, result)
}

func (h *ImageHandler) SoftDelete(c *fiber.Ctx) error {
	imageID, err := parseRequiredUint64(c.Params("imageID"), "imageID")
	if err != nil {
		return httpx.Error(c, err)
	}

	if err := h.service.SoftDelete(c.Context(), imageID); err != nil {
		return httpx.Error(c, err)
	}

	return httpx.Message(c, fiber.StatusOK, "image deleted")
}
