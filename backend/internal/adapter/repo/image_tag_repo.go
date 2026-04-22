package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	mysqlerr "github.com/go-sql-driver/mysql"

	"github.com/example/test-full-stack-developer/backend/internal/apperrors"
	"github.com/example/test-full-stack-developer/backend/internal/port"
)

type ImageTagRepo struct {
	db *sql.DB
}

func NewImageTagRepo(db *sql.DB) port.ImageTagRepo {
	return &ImageTagRepo{db: db}
}

func (r *ImageTagRepo) Attach(ctx context.Context, imageID, tagID uint64) error {
	query := `INSERT INTO image_tags (image_id, tag_id) VALUES (?, ?)`
	if _, err := r.db.ExecContext(ctx, query, imageID, tagID); err != nil {
		var mysqlError *mysqlerr.MySQLError
		if errors.As(err, &mysqlError) && mysqlError.Number == 1062 {
			return apperrors.NewConflict("tag is already attached to image")
		}
		return fmt.Errorf("attach tag to image: %w", err)
	}

	return nil
}

func (r *ImageTagRepo) Detach(ctx context.Context, imageID, tagID uint64) error {
	query := `DELETE FROM image_tags WHERE image_id = ? AND tag_id = ?`
	result, err := r.db.ExecContext(ctx, query, imageID, tagID)
	if err != nil {
		return fmt.Errorf("detach tag from image: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("detach tag from image rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("%w: image tag assignment", apperrors.ErrNotFound)
	}

	return nil
}
