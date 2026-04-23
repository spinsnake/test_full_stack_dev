package repo

import (
	"context"
	"errors"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/example/test-full-stack-developer/backend/internal/apperrors"
	"github.com/example/test-full-stack-developer/backend/internal/entities"
)

func TestImageRepoCreate(t *testing.T) {
	db, mock := newSQLMock(t)
	repo := &ImageRepo{db: db}

	now := time.Now().UTC().Truncate(time.Millisecond)

	mock.ExpectExec(mustPattern(
		`INSERT INTO images (
			image_url,
			thumbnail_url,
			width,
			height,
			alt_text,
			source,
			deleted_at
		) VALUES (?, ?, ?, ?, ?, ?, NULL)`,
	)).
		WithArgs("https://placehold.co/1200x900?text=Created", nil, nil, nil, nil, nil).
		WillReturnResult(sqlmock.NewResult(7, 1))

	mock.ExpectQuery(mustPattern(
		`SELECT
			id,
			image_url,
			thumbnail_url,
			width,
			height,
			alt_text,
			source,
			deleted_at,
			created_at,
			updated_at
		FROM images
		WHERE id = ? AND deleted_at IS NULL`,
	)).
		WithArgs(uint64(7)).
		WillReturnRows(sqlmock.NewRows(imageColumns).AddRow(
			imageValues(7, "https://placehold.co/1200x900?text=Created", nil, nil, nil, nil, nil, nil, now, now)...,
		))

	mock.ExpectQuery(imageHydrateTagsQuery(1)).
		WithArgs(uint64(7)).
		WillReturnRows(sqlmock.NewRows([]string{"image_id", "id", "name", "slug"}))

	image, err := repo.Create(context.Background(), entities.CreateImageInput{
		ImageURL: "https://placehold.co/1200x900?text=Created",
	})
	if err != nil {
		t.Fatalf("create image: %v", err)
	}

	if image.ID != 7 {
		t.Fatalf("expected image id 7, got %d", image.ID)
	}
	if image.ImageURL != "https://placehold.co/1200x900?text=Created" {
		t.Fatalf("unexpected image url: %s", image.ImageURL)
	}
	if len(image.Tags) != 0 {
		t.Fatalf("expected no tags, got %d", len(image.Tags))
	}
}

func TestImageRepoGetByIDNotFound(t *testing.T) {
	db, mock := newSQLMock(t)
	repo := &ImageRepo{db: db}

	mock.ExpectQuery(mustPattern(
		`SELECT
			id,
			image_url,
			thumbnail_url,
			width,
			height,
			alt_text,
			source,
			deleted_at,
			created_at,
			updated_at
		FROM images
		WHERE id = ? AND deleted_at IS NULL`,
	)).
		WithArgs(uint64(42)).
		WillReturnRows(sqlmock.NewRows(imageColumns))

	_, err := repo.GetByID(context.Background(), 42)
	if !errors.Is(err, apperrors.ErrNotFound) {
		t.Fatalf("expected not found error, got %v", err)
	}
}

func TestImageRepoListWithTagAndCursor(t *testing.T) {
	db, mock := newSQLMock(t)
	repo := &ImageRepo{db: db}

	now := time.Now().UTC().Truncate(time.Millisecond)
	filter := entities.ImageListFilter{
		Cursor:  func() *uint64 { value := uint64(10); return &value }(),
		Limit:   2,
		TagSlug: "nature",
	}

	mock.ExpectQuery(mustPattern(
		`SELECT DISTINCT
			i.id,
			i.image_url,
			i.thumbnail_url,
			i.width,
			i.height,
			i.alt_text,
			i.source,
			i.deleted_at,
			i.created_at,
			i.updated_at
		FROM images i
		INNER JOIN image_tags it ON it.image_id = i.id
		INNER JOIN tags t ON t.id = it.tag_id
		WHERE i.deleted_at IS NULL
		  AND t.deleted_at IS NULL
		  AND t.slug = ?
		  AND i.id < ?
		ORDER BY i.id DESC LIMIT ?`,
	)).
		WithArgs("nature", uint64(10), 2).
		WillReturnRows(
			sqlmock.NewRows(imageColumns).
				AddRow(imageValues(9, "https://placehold.co/1200x900?text=9", nil, testInt(1200), testInt(900), testString("Nine"), testString("seed"), nil, now, now)...).
				AddRow(imageValues(8, "https://placehold.co/900x1200?text=8", testString("https://placehold.co/400x300?text=8"), testInt(900), testInt(1200), testString("Eight"), testString("seed"), nil, now, now)...),
		)

	mock.ExpectQuery(imageHydrateTagsQuery(2)).
		WithArgs(uint64(9), uint64(8)).
		WillReturnRows(
			sqlmock.NewRows([]string{"image_id", "id", "name", "slug"}).
				AddRow(int64(8), int64(2), "Portrait", "portrait").
				AddRow(int64(9), int64(1), "Nature", "nature").
				AddRow(int64(9), int64(3), "Travel", "travel"),
		)

	images, err := repo.List(context.Background(), filter)
	if err != nil {
		t.Fatalf("list images: %v", err)
	}

	if len(images) != 2 {
		t.Fatalf("expected 2 images, got %d", len(images))
	}
	if images[0].ID != 9 || len(images[0].Tags) != 2 {
		t.Fatalf("unexpected first image payload: %#v", images[0])
	}
	if images[1].ID != 8 || len(images[1].Tags) != 1 {
		t.Fatalf("unexpected second image payload: %#v", images[1])
	}
	if images[0].Tags[0].Slug != "nature" || images[0].Tags[1].Slug != "travel" {
		t.Fatalf("unexpected first image tags: %#v", images[0].Tags)
	}
}

func TestImageRepoUpdateNotFound(t *testing.T) {
	db, mock := newSQLMock(t)
	repo := &ImageRepo{db: db}

	mock.ExpectExec(mustPattern(
		`UPDATE images SET alt_text = ? WHERE id = ? AND deleted_at IS NULL`,
	)).
		WithArgs("updated caption", uint64(99)).
		WillReturnResult(sqlmock.NewResult(0, 0))

	_, err := repo.Update(context.Background(), 99, entities.UpdateImageInput{
		AltText: testString("updated caption"),
	})
	if !errors.Is(err, apperrors.ErrNotFound) {
		t.Fatalf("expected not found error, got %v", err)
	}
}

func TestImageRepoSoftDeleteNotFound(t *testing.T) {
	db, mock := newSQLMock(t)
	repo := &ImageRepo{db: db}

	mock.ExpectExec(mustPattern(
		`UPDATE images
		SET deleted_at = CURRENT_TIMESTAMP(3)
		WHERE id = ? AND deleted_at IS NULL`,
	)).
		WithArgs(uint64(11)).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := repo.SoftDelete(context.Background(), 11)
	if !errors.Is(err, apperrors.ErrNotFound) {
		t.Fatalf("expected not found error, got %v", err)
	}
}

func TestImageRepoExists(t *testing.T) {
	db, mock := newSQLMock(t)
	repo := &ImageRepo{db: db}

	mock.ExpectQuery(mustPattern(
		`SELECT EXISTS(
			SELECT 1
			FROM images
			WHERE id = ? AND deleted_at IS NULL
		)`,
	)).
		WithArgs(uint64(5)).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	exists, err := repo.Exists(context.Background(), 5)
	if err != nil {
		t.Fatalf("exists: %v", err)
	}
	if !exists {
		t.Fatal("expected image to exist")
	}
}
