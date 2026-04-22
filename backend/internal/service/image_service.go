package service

import (
	"context"
	"strings"

	"github.com/example/test-full-stack-developer/backend/internal/apperrors"
	"github.com/example/test-full-stack-developer/backend/internal/entities"
	"github.com/example/test-full-stack-developer/backend/internal/port"
)

type ImageService struct {
	repository       port.ImageRepo
	defaultPageLimit int
	maxPageLimit     int
}

func NewImageService(repository port.ImageRepo, defaultPageLimit, maxPageLimit int) port.ImageService {
	if defaultPageLimit <= 0 {
		defaultPageLimit = 12
	}
	if maxPageLimit <= 0 {
		maxPageLimit = 60
	}

	return &ImageService{
		repository:       repository,
		defaultPageLimit: defaultPageLimit,
		maxPageLimit:     maxPageLimit,
	}
}

func (s *ImageService) Create(ctx context.Context, input entities.CreateImageInput) (entities.Image, error) {
	input.ImageURL = strings.TrimSpace(input.ImageURL)
	if input.ImageURL == "" {
		return entities.Image{}, apperrors.NewInvalidInput("image_url is required")
	}

	if err := validateDimensions(input.Width, input.Height); err != nil {
		return entities.Image{}, err
	}

	return s.repository.Create(ctx, input)
}

func (s *ImageService) GetByID(ctx context.Context, id uint64) (entities.Image, error) {
	return s.repository.GetByID(ctx, id)
}

func (s *ImageService) List(ctx context.Context, cursor *uint64, limit *int, tagSlug string) ([]entities.Image, error) {
	resolvedLimit := s.defaultPageLimit
	if limit != nil {
		resolvedLimit = *limit
	}
	if resolvedLimit <= 0 {
		resolvedLimit = s.defaultPageLimit
	}
	if resolvedLimit > s.maxPageLimit {
		resolvedLimit = s.maxPageLimit
	}

	filter := entities.ImageListFilter{
		Cursor:  cursor,
		Limit:   resolvedLimit,
		TagSlug: strings.TrimSpace(strings.ToLower(tagSlug)),
	}

	return s.repository.List(ctx, filter)
}

func (s *ImageService) Update(ctx context.Context, id uint64, input entities.UpdateImageInput) (entities.Image, error) {
	if input.ImageURL != nil {
		trimmed := strings.TrimSpace(*input.ImageURL)
		if trimmed == "" {
			return entities.Image{}, apperrors.NewInvalidInput("image_url cannot be empty")
		}
		input.ImageURL = &trimmed
	}

	if err := validateDimensions(input.Width, input.Height); err != nil {
		return entities.Image{}, err
	}

	return s.repository.Update(ctx, id, input)
}

func (s *ImageService) SoftDelete(ctx context.Context, id uint64) error {
	return s.repository.SoftDelete(ctx, id)
}

func validateDimensions(width, height *int) error {
	if width != nil && *width <= 0 {
		return apperrors.NewInvalidInput("width must be greater than 0")
	}
	if height != nil && *height <= 0 {
		return apperrors.NewInvalidInput("height must be greater than 0")
	}
	return nil
}
