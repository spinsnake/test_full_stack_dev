package infra

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	mysqlmigrate "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var defaultMigrationDirs = []string{
	"migration",
	filepath.Join("backend", "migration"),
	"/app/migration",
}

func RunMigrations(db *sql.DB, migrationDir string) error {
	resolvedDir, err := resolveMigrationDir(migrationDir)
	if err != nil {
		return err
	}

	driver, err := mysqlmigrate.WithInstance(db, &mysqlmigrate.Config{})
	if err != nil {
		return fmt.Errorf("create mysql migration driver: %w", err)
	}

	runner, err := migrate.NewWithDatabaseInstance(fileSourceURI(resolvedDir), "mysql", driver)
	if err != nil {
		return fmt.Errorf("create migration runner: %w", err)
	}

	if err := runner.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("apply migrations: %w", err)
	}

	return nil
}

func resolveMigrationDir(custom string) (string, error) {
	candidates := defaultMigrationDirs
	if strings.TrimSpace(custom) != "" {
		candidates = []string{custom}
	}

	for _, candidate := range candidates {
		absolute, err := filepath.Abs(candidate)
		if err != nil {
			continue
		}

		info, err := os.Stat(absolute)
		if err != nil || !info.IsDir() {
			continue
		}

		return absolute, nil
	}

	return "", fmt.Errorf("migration directory not found")
}

func fileSourceURI(path string) string {
	normalized := filepath.ToSlash(path)
	if !strings.HasPrefix(normalized, "/") {
		normalized = "/" + normalized
	}

	return "file://" + normalized
}
