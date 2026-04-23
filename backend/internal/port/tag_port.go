package port

import (
	"context"

	"github.com/example/test-full-stack-developer/backend/internal/entities"
)

type TagRepo interface {
	Create(ctx context.Context, input entities.CreateTagInput) (entities.Tag, error)
	GetByID(ctx context.Context, id uint64) (entities.Tag, error)
	List(ctx context.Context) ([]entities.Tag, error)
	Update(ctx context.Context, id uint64, input entities.UpdateTagInput) (entities.Tag, error)
	SoftDelete(ctx context.Context, id uint64) error
	Exists(ctx context.Context, id uint64) (bool, error)
}

type TagService interface {
	Create(ctx context.Context, input entities.CreateTagInput) (entities.Tag, error)
	GetByID(ctx context.Context, id uint64) (entities.Tag, error)
	List(ctx context.Context) ([]entities.Tag, error)
	Update(ctx context.Context, id uint64, input entities.UpdateTagInput) (entities.Tag, error)
	SoftDelete(ctx context.Context, id uint64) error
}
