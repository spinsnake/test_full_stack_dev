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

const (
	demoSeedUpFile   = "000002_seed_demo_data.up.sql"
	demoSeedDownFile = "000002_seed_demo_data.down.sql"
	demoSeedMarker   = "seed:placehold.co"
)

var defaultMigrationDirs = []string{
	"migration",
	filepath.Join("backend", "migration"),
	"/app/migration",
}

func RunMigrations(db *sql.DB, migrationDir string, includeMockData bool) error {
	resolvedDir, err := resolveMigrationDir(migrationDir)
	if err != nil {
		return err
	}

	schemaDir, cleanup, err := prepareSchemaMigrationDir(resolvedDir)
	if err != nil {
		return err
	}
	defer cleanup()

	driver, err := mysqlmigrate.WithInstance(db, &mysqlmigrate.Config{})
	if err != nil {
		return fmt.Errorf("create mysql migration driver: %w", err)
	}

	runner, err := migrate.NewWithDatabaseInstance(fileSourceURI(schemaDir), "mysql", driver)
	if err != nil {
		return fmt.Errorf("create migration runner: %w", err)
	}

	if err := runner.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("apply migrations: %w", err)
	}

	if includeMockData {
		if err := ensureDemoSeedData(db, resolvedDir); err != nil {
			return err
		}
	}

	return nil
}

func prepareSchemaMigrationDir(sourceDir string) (string, func(), error) {
	tempDir, err := os.MkdirTemp("", "gallery-schema-migrations-*")
	if err != nil {
		return "", nil, fmt.Errorf("create temp migration dir: %w", err)
	}

	cleanup := func() {
		_ = os.RemoveAll(tempDir)
	}

	entries, err := os.ReadDir(sourceDir)
	if err != nil {
		cleanup()
		return "", nil, fmt.Errorf("read migration dir: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		sourcePath := filepath.Join(sourceDir, entry.Name())
		targetPath := filepath.Join(tempDir, entry.Name())

		if isDemoSeedMigration(entry.Name()) {
			if err := os.WriteFile(targetPath, []byte("SELECT 1;\n"), 0o644); err != nil {
				cleanup()
				return "", nil, fmt.Errorf("write stub migration %s: %w", entry.Name(), err)
			}
			continue
		}

		content, err := os.ReadFile(sourcePath)
		if err != nil {
			cleanup()
			return "", nil, fmt.Errorf("read migration file %s: %w", entry.Name(), err)
		}

		if err := os.WriteFile(targetPath, content, 0o644); err != nil {
			cleanup()
			return "", nil, fmt.Errorf("copy migration file %s: %w", entry.Name(), err)
		}
	}

	return tempDir, cleanup, nil
}

func ensureDemoSeedData(db *sql.DB, migrationDir string) error {
	seeded, err := hasDemoSeedData(db)
	if err != nil {
		return err
	}
	if seeded {
		return nil
	}

	seedPath := filepath.Join(migrationDir, demoSeedUpFile)
	body, err := os.ReadFile(seedPath)
	if err != nil {
		return fmt.Errorf("read demo seed file: %w", err)
	}

	if err := execSQLStatements(db, string(body)); err != nil {
		return fmt.Errorf("apply demo seed file: %w", err)
	}

	return nil
}

func hasDemoSeedData(db *sql.DB) (bool, error) {
	const query = `
		SELECT EXISTS(
			SELECT 1
			FROM images
			WHERE source = ?
			LIMIT 1
		)
	`

	var exists bool
	if err := db.QueryRow(query, demoSeedMarker).Scan(&exists); err != nil {
		return false, fmt.Errorf("check demo seed data: %w", err)
	}

	return exists, nil
}

func execSQLStatements(db *sql.DB, script string) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("begin seed transaction: %w", err)
	}

	for _, statement := range splitSQLStatements(script) {
		if _, err := tx.Exec(statement); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("exec seed statement: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit seed transaction: %w", err)
	}

	return nil
}

func splitSQLStatements(script string) []string {
	parts := strings.Split(script, ";")
	statements := make([]string, 0, len(parts))

	for _, part := range parts {
		statement := strings.TrimSpace(part)
		if statement == "" {
			continue
		}

		lines := strings.Split(statement, "\n")
		cleanLines := make([]string, 0, len(lines))
		for _, line := range lines {
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, "--") {
				continue
			}
			cleanLines = append(cleanLines, line)
		}

		cleaned := strings.TrimSpace(strings.Join(cleanLines, "\n"))
		if cleaned == "" {
			continue
		}

		statements = append(statements, cleaned)
	}

	return statements
}

func isDemoSeedMigration(name string) bool {
	return name == demoSeedUpFile || name == demoSeedDownFile
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
