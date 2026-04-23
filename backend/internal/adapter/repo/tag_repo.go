package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	mysqlerr "github.com/go-sql-driver/mysql"

	"github.com/example/test-full-stack-developer/backend/internal/apperrors"
	"github.com/example/test-full-stack-developer/backend/internal/entities"
	"github.com/example/test-full-stack-developer/backend/internal/port"
)

type TagRepo struct {
	db *sql.DB
}

func NewTagRepo(db *sql.DB) port.TagRepo {
	return &TagRepo{db: db}
}

func (r *TagRepo) Create(ctx context.Context, input entities.CreateTagInput) (entities.Tag, error) {
	query := `
		INSERT INTO tags (
			name,
			slug,
			deleted_at
		) VALUES (?, ?, NULL)
	`

	result, err := r.db.ExecContext(ctx, query, input.Name, input.Slug)
	if err != nil {
		return entities.Tag{}, normalizeTagWriteError(err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return entities.Tag{}, fmt.Errorf("tag last insert id: %w", err)
	}

	return r.GetByID(ctx, uint64(id))
}

func (r *TagRepo) GetByID(ctx context.Context, id uint64) (entities.Tag, error) {
	query := `
		SELECT
			id,
			name,
			slug,
			deleted_at,
			created_at,
			updated_at
		FROM tags
		WHERE id = ? AND deleted_at IS NULL
	`

	row := r.db.QueryRowContext(ctx, query, id)
	tag, err := scanTag(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return entities.Tag{}, fmt.Errorf("%w: tag %d", apperrors.ErrNotFound, id)
		}
		return entities.Tag{}, fmt.Errorf("get tag: %w", err)
	}

	return tag, nil
}

func (r *TagRepo) List(ctx context.Context) ([]entities.Tag, error) {
	query := `
		SELECT
			id,
			name,
			slug,
			deleted_at,
			created_at,
			updated_at
		FROM tags
		WHERE deleted_at IS NULL
		ORDER BY name ASC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list tags: %w", err)
	}
	defer rows.Close()

	tags := make([]entities.Tag, 0)
	for rows.Next() {
		tag, err := scanTag(rows)
		if err != nil {
			return nil, fmt.Errorf("scan tag: %w", err)
		}
		tags = append(tags, tag)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate tags: %w", err)
	}

	return tags, nil
}

func (r *TagRepo) Update(ctx context.Context, id uint64, input entities.UpdateTagInput) (entities.Tag, error) {
	setParts := make([]string, 0, 2)
	args := make([]any, 0, 3)

	if input.Name != nil {
		setParts = append(setParts, "name = ?")
		args = append(args, *input.Name)
	}
	if input.Slug != nil {
		setParts = append(setParts, "slug = ?")
		args = append(args, *input.Slug)
	}

	if len(setParts) == 0 {
		return r.GetByID(ctx, id)
	}

	args = append(args, id)
	query := fmt.Sprintf(
		"UPDATE tags SET %s WHERE id = ? AND deleted_at IS NULL",
		strings.Join(setParts, ", "),
	)

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return entities.Tag{}, normalizeTagWriteError(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return entities.Tag{}, fmt.Errorf("update tag rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return entities.Tag{}, fmt.Errorf("%w: tag %d", apperrors.ErrNotFound, id)
	}

	return r.GetByID(ctx, id)
}

func (r *TagRepo) SoftDelete(ctx context.Context, id uint64) error {
	query := `
		UPDATE tags
		SET deleted_at = CURRENT_TIMESTAMP(3)
		WHERE id = ? AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("soft delete tag: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("soft delete tag rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("%w: tag %d", apperrors.ErrNotFound, id)
	}

	return nil
}

func (r *TagRepo) Exists(ctx context.Context, id uint64) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1
			FROM tags
			WHERE id = ? AND deleted_at IS NULL
		)
	`

	var exists bool
	if err := r.db.QueryRowContext(ctx, query, id).Scan(&exists); err != nil {
		return false, fmt.Errorf("check tag exists: %w", err)
	}

	return exists, nil
}

func normalizeTagWriteError(err error) error {
	var mysqlError *mysqlerr.MySQLError
	if errors.As(err, &mysqlError) && mysqlError.Number == 1062 {
		return apperrors.NewConflict("tag name or slug already exists")
	}
	return fmt.Errorf("write tag: %w", err)
}

func scanTag(scanner interface{ Scan(dest ...any) error }) (entities.Tag, error) {
	var (
		tag       entities.Tag
		deletedAt sql.NullTime
		createdAt time.Time
		updatedAt time.Time
	)

	if err := scanner.Scan(
		&tag.ID,
		&tag.Name,
		&tag.Slug,
		&deletedAt,
		&createdAt,
		&updatedAt,
	); err != nil {
		return entities.Tag{}, err
	}

	if deletedAt.Valid {
		value := deletedAt.Time
		tag.DeletedAt = &value
	}
	tag.CreatedAt = createdAt
	tag.UpdatedAt = updatedAt

	return tag, nil
}
