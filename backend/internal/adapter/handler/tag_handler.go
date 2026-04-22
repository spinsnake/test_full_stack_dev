package handler

import (
	"github.com/gofiber/fiber/v2"

	"github.com/example/test-full-stack-developer/backend/internal/entities"
	"github.com/example/test-full-stack-developer/backend/internal/httpx"
	"github.com/example/test-full-stack-developer/backend/internal/port"
)

type TagHandler struct {
	service port.TagService
}

func NewTagHandler(service port.TagService) *TagHandler {
	return &TagHandler{service: service}
}

func (h *TagHandler) Create(c *fiber.Ctx) error {
	var input entities.CreateTagInput
	if err := c.BodyParser(&input); err != nil {
		return httpx.Error(c, err)
	}

	result, err := h.service.Create(c.Context(), input)
	if err != nil {
		return httpx.Error(c, err)
	}

	return httpx.Data(c, fiber.StatusCreated, result)
}

func (h *TagHandler) List(c *fiber.Ctx) error {
	result, err := h.service.List(c.Context())
	if err != nil {
		return httpx.Error(c, err)
	}

	return httpx.Data(c, fiber.StatusOK, result)
}

func (h *TagHandler) GetByID(c *fiber.Ctx) error {
	tagID, err := parseRequiredUint64(c.Params("tagID"), "tagID")
	if err != nil {
		return httpx.Error(c, err)
	}

	result, err := h.service.GetByID(c.Context(), tagID)
	if err != nil {
		return httpx.Error(c, err)
	}

	return httpx.Data(c, fiber.StatusOK, result)
}

func (h *TagHandler) Update(c *fiber.Ctx) error {
	tagID, err := parseRequiredUint64(c.Params("tagID"), "tagID")
	if err != nil {
		return httpx.Error(c, err)
	}

	var input entities.UpdateTagInput
	if err := c.BodyParser(&input); err != nil {
		return httpx.Error(c, err)
	}

	result, err := h.service.Update(c.Context(), tagID, input)
	if err != nil {
		return httpx.Error(c, err)
	}

	return httpx.Data(c, fiber.StatusOK, result)
}

func (h *TagHandler) SoftDelete(c *fiber.Ctx) error {
	tagID, err := parseRequiredUint64(c.Params("tagID"), "tagID")
	if err != nil {
		return httpx.Error(c, err)
	}

	if err := h.service.SoftDelete(c.Context(), tagID); err != nil {
		return httpx.Error(c, err)
	}

	return httpx.Message(c, fiber.StatusOK, "tag deleted")
}
