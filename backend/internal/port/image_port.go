package port

import (
	"context"

	"github.com/example/test-full-stack-developer/backend/internal/entities"
)

type ImageRepo interface {
	Create(ctx context.Context, input entities.CreateImageInput) (entities.Image, error)
	GetByID(ctx context.Context, id uint64) (entities.Image, error)
	List(ctx context.Context, filter entities.ImageListFilter) ([]entities.Image, error)
	Update(ctx context.Context, id uint64, input entities.UpdateImageInput) (entities.Image, error)
	SoftDelete(ctx context.Context, id uint64) error
	Exists(ctx context.Context, id uint64) (bool, error)
}

type ImageService interface {
	Create(ctx context.Context, input entities.CreateImageInput) (entities.Image, error)
	GetByID(ctx context.Context, id uint64) (entities.Image, error)
	List(ctx context.Context, cursor *uint64, limit *int, tagSlug string) ([]entities.Image, error)
	Update(ctx context.Context, id uint64, input entities.UpdateImageInput) (entities.Image, error)
	SoftDelete(ctx context.Context, id uint64) error
}
