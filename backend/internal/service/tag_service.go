package service

import (
	"context"
	"strings"

	"github.com/example/test-full-stack-developer/backend/internal/apperrors"
	"github.com/example/test-full-stack-developer/backend/internal/entities"
	"github.com/example/test-full-stack-developer/backend/internal/port"
)

type TagService struct {
	repository port.TagRepo
}

func NewTagService(repository port.TagRepo) port.TagService {
	return &TagService{repository: repository}
}

func (s *TagService) Create(ctx context.Context, input entities.CreateTagInput) (entities.Tag, error) {
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return entities.Tag{}, apperrors.NewInvalidInput("name is required")
	}

	slug := ""
	if input.Slug != nil {
		slug = *input.Slug
	}
	slug = normalizeSlug(slug, name)
	if slug == "" {
		return entities.Tag{}, apperrors.NewInvalidInput("slug is invalid")
	}

	return s.repository.Create(ctx, entities.CreateTagInput{
		Name: name,
		Slug: &slug,
	})
}

func (s *TagService) GetByID(ctx context.Context, id uint64) (entities.Tag, error) {
	return s.repository.GetByID(ctx, id)
}

func (s *TagService) List(ctx context.Context) ([]entities.Tag, error) {
	return s.repository.List(ctx)
}

func (s *TagService) Update(ctx context.Context, id uint64, input entities.UpdateTagInput) (entities.Tag, error) {
	if input.Name != nil {
		name := strings.TrimSpace(*input.Name)
		if name == "" {
			return entities.Tag{}, apperrors.NewInvalidInput("name cannot be empty")
		}
		input.Name = &name
	}

	if input.Slug != nil {
		var name string
		if input.Name != nil {
			name = *input.Name
		}
		slug := normalizeSlug(*input.Slug, name)
		if slug == "" {
			return entities.Tag{}, apperrors.NewInvalidInput("slug is invalid")
		}
		input.Slug = &slug
	}

	return s.repository.Update(ctx, id, input)
}

func (s *TagService) SoftDelete(ctx context.Context, id uint64) error {
	return s.repository.SoftDelete(ctx, id)
}

func normalizeSlug(rawSlug, fallback string) string {
	source := strings.TrimSpace(strings.ToLower(rawSlug))
	if source == "" {
		source = strings.TrimSpace(strings.ToLower(fallback))
	}
	if source == "" {
		return ""
	}

	var builder strings.Builder
	lastWasDash := false
	for _, char := range source {
		switch {
		case char >= 'a' && char <= 'z':
			builder.WriteRune(char)
			lastWasDash = false
		case char >= '0' && char <= '9':
			builder.WriteRune(char)
			lastWasDash = false
		case char == ' ' || char == '-' || char == '_' || char == '.':
			if !lastWasDash && builder.Len() > 0 {
				builder.WriteByte('-')
				lastWasDash = true
			}
		}
	}

	return strings.Trim(builder.String(), "-")
}
