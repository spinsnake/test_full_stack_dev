package repo

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/example/test-full-stack-developer/backend/internal/apperrors"
	"github.com/example/test-full-stack-developer/backend/internal/entities"
	sqlmock "github.com/DATA-DOG/go-sqlmock"
	mysqlerr "github.com/go-sql-driver/mysql"
)

func TestTagRepoCreateConflict(t *testing.T) {
	db, mock := newSQLMock(t)
	repo := &TagRepo{db: db}

	mock.ExpectExec(mustPattern(
		`INSERT INTO tags (
			name,
			slug,
			is_active,
			deleted_at
		) VALUES (?, ?, 1, NULL)`,
	)).
		WithArgs("Travel", "travel").
		WillReturnError(&mysqlerr.MySQLError{Number: 1062, Message: "Duplicate entry"})

	_, err := repo.Create(context.Background(), entities.CreateTagInput{
		Name: "Travel",
		Slug: testString("travel"),
	})
	if !errors.Is(err, apperrors.ErrConflict) {
		t.Fatalf("expected conflict error, got %v", err)
	}
}

func TestTagRepoUpdate(t *testing.T) {
	db, mock := newSQLMock(t)
	repo := &TagRepo{db: db}

	now := time.Now().UTC().Truncate(time.Millisecond)

	mock.ExpectExec(mustPattern(
		`UPDATE tags SET name = ?, slug = ? WHERE id = ? AND deleted_at IS NULL AND is_active = 1`,
	)).
		WithArgs("City", "city", uint64(3)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectQuery(mustPattern(
		`SELECT
			id,
			name,
			slug,
			is_active,
			deleted_at,
			created_at,
			updated_at
		FROM tags
		WHERE id = ? AND deleted_at IS NULL AND is_active = 1`,
	)).
		WithArgs(uint64(3)).
		WillReturnRows(sqlmock.NewRows(tagColumns).AddRow(
			tagValues(3, "City", "city", true, nil, now, now)...,
		))

	tag, err := repo.Update(context.Background(), 3, entities.UpdateTagInput{
		Name: testString("City"),
		Slug: testString("city"),
	})
	if err != nil {
		t.Fatalf("update tag: %v", err)
	}

	if tag.Name != "City" || tag.Slug != "city" {
		t.Fatalf("unexpected tag payload: %#v", tag)
	}
}

func TestTagRepoList(t *testing.T) {
	db, mock := newSQLMock(t)
	repo := &TagRepo{db: db}

	now := time.Now().UTC().Truncate(time.Millisecond)

	mock.ExpectQuery(mustPattern(
		`SELECT
			id,
			name,
			slug,
			is_active,
			deleted_at,
			created_at,
			updated_at
		FROM tags
		WHERE deleted_at IS NULL AND is_active = 1
		ORDER BY name ASC`,
	)).
		WillReturnRows(
			sqlmock.NewRows(tagColumns).
				AddRow(tagValues(1, "City", "city", true, nil, now, now)...).
				AddRow(tagValues(2, "Travel", "travel", true, nil, now, now)...),
		)

	tags, err := repo.List(context.Background())
	if err != nil {
		t.Fatalf("list tags: %v", err)
	}

	if len(tags) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(tags))
	}
	if tags[0].Name != "City" || tags[1].Slug != "travel" {
		t.Fatalf("unexpected tag list: %#v", tags)
	}
}

func TestTagRepoSoftDeleteNotFound(t *testing.T) {
	db, mock := newSQLMock(t)
	repo := &TagRepo{db: db}

	mock.ExpectExec(mustPattern(
		`UPDATE tags
		SET is_active = 0,
		    deleted_at = CURRENT_TIMESTAMP(3)
		WHERE id = ? AND deleted_at IS NULL AND is_active = 1`,
	)).
		WithArgs(uint64(8)).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := repo.SoftDelete(context.Background(), 8)
	if !errors.Is(err, apperrors.ErrNotFound) {
		t.Fatalf("expected not found error, got %v", err)
	}
}

func TestTagRepoExistsActive(t *testing.T) {
	db, mock := newSQLMock(t)
	repo := &TagRepo{db: db}

	mock.ExpectQuery(mustPattern(
		`SELECT EXISTS(
			SELECT 1
			FROM tags
			WHERE id = ? AND deleted_at IS NULL AND is_active = 1
		)`,
	)).
		WithArgs(uint64(4)).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	exists, err := repo.ExistsActive(context.Background(), 4)
	if err != nil {
		t.Fatalf("exists active: %v", err)
	}
	if !exists {
		t.Fatal("expected tag to exist")
	}
}
