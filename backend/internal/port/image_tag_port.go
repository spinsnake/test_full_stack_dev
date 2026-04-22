package port

import (
	"context"

	"github.com/example/test-full-stack-developer/backend/internal/entities"
)

type ImageTagRepo interface {
	Attach(ctx context.Context, imageID, tagID uint64) error
	Detach(ctx context.Context, imageID, tagID uint64) error
}

type ImageTagService interface {
	Attach(ctx context.Context, imageID uint64, input entities.AttachTagToImageInput) (entities.ImageTagAssignment, error)
	Detach(ctx context.Context, imageID, tagID uint64) error
}
