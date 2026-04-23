package service

import (
	"context"
	"fmt"

	"github.com/example/test-full-stack-developer/backend/internal/apperrors"
	"github.com/example/test-full-stack-developer/backend/internal/entities"
	"github.com/example/test-full-stack-developer/backend/internal/port"
)

type ImageTagService struct {
	repository port.ImageTagRepo
	imageRepo  port.ImageRepo
	tagRepo    port.TagRepo
}

func NewImageTagService(repository port.ImageTagRepo, imageRepo port.ImageRepo, tagRepo port.TagRepo) port.ImageTagService {
	return &ImageTagService{
		repository: repository,
		imageRepo:  imageRepo,
		tagRepo:    tagRepo,
	}
}

func (s *ImageTagService) Attach(ctx context.Context, imageID uint64, input entities.AttachTagToImageInput) (entities.ImageTagAssignment, error) {
	if input.TagID == 0 {
		return entities.ImageTagAssignment{}, apperrors.NewInvalidInput("tag_id is required")
	}

	imageExists, err := s.imageRepo.Exists(ctx, imageID)
	if err != nil {
		return entities.ImageTagAssignment{}, err
	}
	if !imageExists {
		return entities.ImageTagAssignment{}, fmt.Errorf("%w: image %d", apperrors.ErrNotFound, imageID)
	}

	tagExists, err := s.tagRepo.Exists(ctx, input.TagID)
	if err != nil {
		return entities.ImageTagAssignment{}, err
	}
	if !tagExists {
		return entities.ImageTagAssignment{}, fmt.Errorf("%w: tag %d", apperrors.ErrNotFound, input.TagID)
	}

	if err := s.repository.Attach(ctx, imageID, input.TagID); err != nil {
		return entities.ImageTagAssignment{}, err
	}

	return entities.ImageTagAssignment{
		ImageID: imageID,
		TagID:   input.TagID,
	}, nil
}

func (s *ImageTagService) Detach(ctx context.Context, imageID, tagID uint64) error {
	return s.repository.Detach(ctx, imageID, tagID)
}
