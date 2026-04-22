package handler

import (
	"github.com/gofiber/fiber/v2"

	"github.com/example/test-full-stack-developer/backend/internal/entities"
	"github.com/example/test-full-stack-developer/backend/internal/httpx"
	"github.com/example/test-full-stack-developer/backend/internal/port"
)

type ImageTagHandler struct {
	service port.ImageTagService
}

func NewImageTagHandler(service port.ImageTagService) *ImageTagHandler {
	return &ImageTagHandler{service: service}
}

func (h *ImageTagHandler) Attach(c *fiber.Ctx) error {
	imageID, err := parseRequiredUint64(c.Params("imageID"), "imageID")
	if err != nil {
		return httpx.Error(c, err)
	}

	var payload entities.AttachTagToImageInput
	if err := c.BodyParser(&payload); err != nil {
		return httpx.Error(c, err)
	}

	result, err := h.service.Attach(c.Context(), imageID, payload)
	if err != nil {
		return httpx.Error(c, err)
	}

	return httpx.Data(c, fiber.StatusCreated, result)
}

func (h *ImageTagHandler) Detach(c *fiber.Ctx) error {
	imageID, err := parseRequiredUint64(c.Params("imageID"), "imageID")
	if err != nil {
		return httpx.Error(c, err)
	}

	tagID, err := parseRequiredUint64(c.Params("tagID"), "tagID")
	if err != nil {
		return httpx.Error(c, err)
	}

	if err := h.service.Detach(c.Context(), imageID, tagID); err != nil {
		return httpx.Error(c, err)
	}

	return httpx.Message(c, fiber.StatusOK, "tag removed from image")
}
