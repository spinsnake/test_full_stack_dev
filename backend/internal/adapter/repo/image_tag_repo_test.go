package repo

import (
	"context"
	"errors"
	"testing"

	"github.com/example/test-full-stack-developer/backend/internal/apperrors"
	sqlmock "github.com/DATA-DOG/go-sqlmock"
	mysqlerr "github.com/go-sql-driver/mysql"
)

func TestImageTagRepoAttach(t *testing.T) {
	db, mock := newSQLMock(t)
	repo := &ImageTagRepo{db: db}

	mock.ExpectExec(mustPattern(`INSERT INTO image_tags (image_id, tag_id) VALUES (?, ?)`)).
		WithArgs(uint64(7), uint64(9)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err := repo.Attach(context.Background(), 7, 9); err != nil {
		t.Fatalf("attach tag to image: %v", err)
	}
}

func TestImageTagRepoAttachConflict(t *testing.T) {
	db, mock := newSQLMock(t)
	repo := &ImageTagRepo{db: db}

	mock.ExpectExec(mustPattern(`INSERT INTO image_tags (image_id, tag_id) VALUES (?, ?)`)).
		WithArgs(uint64(7), uint64(9)).
		WillReturnError(&mysqlerr.MySQLError{Number: 1062, Message: "Duplicate entry"})

	err := repo.Attach(context.Background(), 7, 9)
	if !errors.Is(err, apperrors.ErrConflict) {
		t.Fatalf("expected conflict error, got %v", err)
	}
}

func TestImageTagRepoDetachNotFound(t *testing.T) {
	db, mock := newSQLMock(t)
	repo := &ImageTagRepo{db: db}

	mock.ExpectExec(mustPattern(`DELETE FROM image_tags WHERE image_id = ? AND tag_id = ?`)).
		WithArgs(uint64(7), uint64(9)).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := repo.Detach(context.Background(), 7, 9)
	if !errors.Is(err, apperrors.ErrNotFound) {
		t.Fatalf("expected not found error, got %v", err)
	}
}
