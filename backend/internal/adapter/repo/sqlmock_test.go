package repo

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
)

var (
	imageColumns = []string{
		"id",
		"image_url",
		"thumbnail_url",
		"width",
		"height",
		"alt_text",
		"source",
		"is_active",
		"deleted_at",
		"created_at",
		"updated_at",
	}
	tagColumns = []string{
		"id",
		"name",
		"slug",
		"is_active",
		"deleted_at",
		"created_at",
		"updated_at",
	}
)

func newSQLMock(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	t.Helper()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	if err != nil {
		t.Fatalf("create sqlmock: %v", err)
	}

	t.Cleanup(func() {
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet sqlmock expectations: %v", err)
		}
		_ = db.Close()
	})

	return db, mock
}

func mustPattern(query string) string {
	return regexp.QuoteMeta(strings.Join(strings.Fields(query), " "))
}

func imageHydrateTagsQuery(imageCount int) string {
	placeholders := strings.TrimSuffix(strings.Repeat("?, ", imageCount), ", ")
	return mustPattern(fmt.Sprintf(
		`SELECT it.image_id, t.id, t.name, t.slug
		FROM image_tags it
		INNER JOIN tags t ON t.id = it.tag_id
		WHERE it.image_id IN (%s)
		  AND t.is_active = 1
		  AND t.deleted_at IS NULL
		ORDER BY it.image_id, t.name`,
		placeholders,
	))
}

func imageValues(
	id int64,
	imageURL string,
	thumbnail *string,
	width *int,
	height *int,
	altText *string,
	source *string,
	isActive bool,
	deletedAt *time.Time,
	createdAt time.Time,
	updatedAt time.Time,
) []driver.Value {
	return []driver.Value{
		id,
		imageURL,
		nilString(thumbnail),
		nilInt(width),
		nilInt(height),
		nilString(altText),
		nilString(source),
		isActive,
		nilTime(deletedAt),
		createdAt,
		updatedAt,
	}
}

func tagValues(
	id int64,
	name string,
	slug string,
	isActive bool,
	deletedAt *time.Time,
	createdAt time.Time,
	updatedAt time.Time,
) []driver.Value {
	return []driver.Value{
		id,
		name,
		slug,
		isActive,
		nilTime(deletedAt),
		createdAt,
		updatedAt,
	}
}

func nilString(value *string) any {
	if value == nil {
		return nil
	}
	return *value
}

func nilInt(value *int) any {
	if value == nil {
		return nil
	}
	return int64(*value)
}

func nilTime(value *time.Time) any {
	if value == nil {
		return nil
	}
	return *value
}

func testString(value string) *string {
	return &value
}

func testInt(value int) *int {
	return &value
}
