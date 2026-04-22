package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func parseRequiredUint64(value, field string) (uint64, error) {
	parsed, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0, fiber.NewError(fiber.StatusBadRequest, field+" must be a valid unsigned integer")
	}

	return parsed, nil
}

func parseOptionalUint64(value string) (*uint64, error) {
	if value == "" {
		return nil, nil
	}

	parsed, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, "cursor must be a valid unsigned integer")
	}

	return &parsed, nil
}

func parseOptionalInt(value string) (*int, error) {
	if value == "" {
		return nil, nil
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, "limit must be a valid integer")
	}

	return &parsed, nil
}
