package infra

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	schemaUpFile   = "000001_create_gallery_tables.up.sql"
	demoSeedUpFile = "000002_seed_demo_data.up.sql"
	demoSeedMarker = "seed:placehold.co"
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

	if err := applySQLFile(db, filepath.Join(resolvedDir, schemaUpFile)); err != nil {
		return fmt.Errorf("apply schema file: %w", err)
	}

	if includeMockData {
		if err := ensureDemoSeedData(db, resolvedDir); err != nil {
			return err
		}
	}

	return nil
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
	if err := applySQLFile(db, seedPath); err != nil {
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

func applySQLFile(db *sql.DB, path string) error {
	body, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read sql file %s: %w", filepath.Base(path), err)
	}

	return execSQLStatements(db, string(body))
}

func execSQLStatements(db *sql.DB, script string) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("begin sql transaction: %w", err)
	}

	for _, statement := range splitSQLStatements(script) {
		if _, err := tx.Exec(statement); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("exec sql statement: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit sql transaction: %w", err)
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
