package repo

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/example/test-full-stack-developer/backend/internal/apperrors"
	"github.com/example/test-full-stack-developer/backend/internal/entities"
	"github.com/example/test-full-stack-developer/backend/internal/port"
)

type ImageRepo struct {
	db *sql.DB
}

func NewImageRepo(db *sql.DB) port.ImageRepo {
	return &ImageRepo{db: db}
}

func (r *ImageRepo) Create(ctx context.Context, input entities.CreateImageInput) (entities.Image, error) {
	query := `
		INSERT INTO images (
			image_url,
			thumbnail_url,
			width,
			height,
			alt_text,
			source,
			is_active,
			deleted_at
		) VALUES (?, ?, ?, ?, ?, ?, 1, NULL)
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		input.ImageURL,
		input.ThumbnailURL,
		input.Width,
		input.Height,
		input.AltText,
		input.Source,
	)
	if err != nil {
		return entities.Image{}, fmt.Errorf("create image: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return entities.Image{}, fmt.Errorf("image last insert id: %w", err)
	}

	return r.GetByID(ctx, uint64(id))
}

func (r *ImageRepo) GetByID(ctx context.Context, id uint64) (entities.Image, error) {
	query := `
		SELECT
			id,
			image_url,
			thumbnail_url,
			width,
			height,
			alt_text,
			source,
			is_active,
			deleted_at,
			created_at,
			updated_at
		FROM images
		WHERE id = ? AND deleted_at IS NULL AND is_active = 1
	`

	row := r.db.QueryRowContext(ctx, query, id)
	image, err := scanImage(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return entities.Image{}, fmt.Errorf("%w: image %d", apperrors.ErrNotFound, id)
		}
		return entities.Image{}, fmt.Errorf("get image: %w", err)
	}

	images := []entities.Image{image}
	if err := r.hydrateTags(ctx, images); err != nil {
		return entities.Image{}, err
	}

	return images[0], nil
}

func (r *ImageRepo) List(ctx context.Context, filter entities.ImageListFilter) ([]entities.Image, error) {
	args := make([]any, 0, 4)
	var queryBuilder strings.Builder

	queryBuilder.WriteString(`
		SELECT DISTINCT
			i.id,
			i.image_url,
			i.thumbnail_url,
			i.width,
			i.height,
			i.alt_text,
			i.source,
			i.is_active,
			i.deleted_at,
			i.created_at,
			i.updated_at
		FROM images i
	`)

	if filter.TagSlug != "" {
		queryBuilder.WriteString(`
			INNER JOIN image_tags it ON it.image_id = i.id
			INNER JOIN tags t ON t.id = it.tag_id
		`)
	}

	queryBuilder.WriteString(`
		WHERE i.is_active = 1
		  AND i.deleted_at IS NULL
	`)

	if filter.TagSlug != "" {
		queryBuilder.WriteString(`
		  AND t.is_active = 1
		  AND t.deleted_at IS NULL
		  AND t.slug = ?
		`)
		args = append(args, filter.TagSlug)
	}

	if filter.Cursor != nil {
		queryBuilder.WriteString(` AND i.id < ?`)
		args = append(args, *filter.Cursor)
	}

	queryBuilder.WriteString(` ORDER BY i.id DESC LIMIT ?`)
	args = append(args, filter.Limit)

	rows, err := r.db.QueryContext(ctx, queryBuilder.String(), args...)
	if err != nil {
		return nil, fmt.Errorf("list images: %w", err)
	}
	defer rows.Close()

	images := make([]entities.Image, 0, filter.Limit)
	for rows.Next() {
		image, err := scanImage(rows)
		if err != nil {
			return nil, fmt.Errorf("scan image: %w", err)
		}
		images = append(images, image)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate images: %w", err)
	}

	if len(images) == 0 {
		return images, nil
	}

	if err := r.hydrateTags(ctx, images); err != nil {
		return nil, err
	}

	return images, nil
}

func (r *ImageRepo) Update(ctx context.Context, id uint64, input entities.UpdateImageInput) (entities.Image, error) {
	setParts := make([]string, 0, 6)
	args := make([]any, 0, 7)

	if input.ImageURL != nil {
		setParts = append(setParts, "image_url = ?")
		args = append(args, *input.ImageURL)
	}
	if input.ThumbnailURL != nil {
		setParts = append(setParts, "thumbnail_url = ?")
		args = append(args, *input.ThumbnailURL)
	}
	if input.Width != nil {
		setParts = append(setParts, "width = ?")
		args = append(args, *input.Width)
	}
	if input.Height != nil {
		setParts = append(setParts, "height = ?")
		args = append(args, *input.Height)
	}
	if input.AltText != nil {
		setParts = append(setParts, "alt_text = ?")
		args = append(args, *input.AltText)
	}
	if input.Source != nil {
		setParts = append(setParts, "source = ?")
		args = append(args, *input.Source)
	}

	if len(setParts) == 0 {
		return r.GetByID(ctx, id)
	}

	args = append(args, id)
	query := fmt.Sprintf(
		"UPDATE images SET %s WHERE id = ? AND deleted_at IS NULL AND is_active = 1",
		strings.Join(setParts, ", "),
	)

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return entities.Image{}, fmt.Errorf("update image: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return entities.Image{}, fmt.Errorf("update image rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return entities.Image{}, fmt.Errorf("%w: image %d", apperrors.ErrNotFound, id)
	}

	return r.GetByID(ctx, id)
}

func (r *ImageRepo) SoftDelete(ctx context.Context, id uint64) error {
	query := `
		UPDATE images
		SET is_active = 0,
		    deleted_at = CURRENT_TIMESTAMP(3)
		WHERE id = ? AND deleted_at IS NULL AND is_active = 1
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("soft delete image: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("soft delete image rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("%w: image %d", apperrors.ErrNotFound, id)
	}

	return nil
}

func (r *ImageRepo) ExistsActive(ctx context.Context, id uint64) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1
			FROM images
			WHERE id = ? AND deleted_at IS NULL AND is_active = 1
		)
	`

	var exists bool
	if err := r.db.QueryRowContext(ctx, query, id).Scan(&exists); err != nil {
		return false, fmt.Errorf("check image exists: %w", err)
	}

	return exists, nil
}

func (r *ImageRepo) hydrateTags(ctx context.Context, images []entities.Image) error {
	if len(images) == 0 {
		return nil
	}

	placeholders := make([]string, 0, len(images))
	args := make([]any, 0, len(images))
	imageMap := make(map[uint64]*entities.Image, len(images))

	for index := range images {
		placeholders = append(placeholders, "?")
		args = append(args, images[index].ID)
		imageMap[images[index].ID] = &images[index]
		images[index].Tags = []entities.TagSummary{}
	}

	query := fmt.Sprintf(`
		SELECT
			it.image_id,
			t.id,
			t.name,
			t.slug
		FROM image_tags it
		INNER JOIN tags t ON t.id = it.tag_id
		WHERE it.image_id IN (%s)
		  AND t.is_active = 1
		  AND t.deleted_at IS NULL
		ORDER BY it.image_id, t.name
	`, strings.Join(placeholders, ", "))

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("load image tags: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var imageID uint64
		var tag entities.TagSummary
		if err := rows.Scan(&imageID, &tag.ID, &tag.Name, &tag.Slug); err != nil {
			return fmt.Errorf("scan image tag: %w", err)
		}

		image, ok := imageMap[imageID]
		if !ok {
			continue
		}
		image.Tags = append(image.Tags, tag)
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate image tags: %w", err)
	}

	return nil
}

func scanImage(scanner interface{ Scan(dest ...any) error }) (entities.Image, error) {
	var (
		image     entities.Image
		thumbnail sql.NullString
		width     sql.NullInt64
		height    sql.NullInt64
		altText   sql.NullString
		source    sql.NullString
		deletedAt sql.NullTime
		createdAt time.Time
		updatedAt time.Time
	)

	if err := scanner.Scan(
		&image.ID,
		&image.ImageURL,
		&thumbnail,
		&width,
		&height,
		&altText,
		&source,
		&image.IsActive,
		&deletedAt,
		&createdAt,
		&updatedAt,
	); err != nil {
		return entities.Image{}, err
	}

	if thumbnail.Valid {
		value := thumbnail.String
		image.ThumbnailURL = &value
	}
	if width.Valid {
		value := int(width.Int64)
		image.Width = &value
	}
	if height.Valid {
		value := int(height.Int64)
		image.Height = &value
	}
	if altText.Valid {
		value := altText.String
		image.AltText = &value
	}
	if source.Valid {
		value := source.String
		image.Source = &value
	}
	if deletedAt.Valid {
		value := deletedAt.Time
		image.DeletedAt = &value
	}

	image.CreatedAt = createdAt
	image.UpdatedAt = updatedAt

	return image, nil
}
