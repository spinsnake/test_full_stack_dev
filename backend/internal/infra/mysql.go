package infra

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/example/test-full-stack-developer/backend/internal/config"
)

func OpenMySQL(cfg config.Config) (*sql.DB, error) {
	db, err := sql.Open("mysql", cfg.DatabaseDSN)
	if err != nil {
		return nil, fmt.Errorf("sql open: %w", err)
	}

	db.SetMaxOpenConns(cfg.DBMaxOpenConns)
	db.SetMaxIdleConns(cfg.DBMaxIdleConns)
	db.SetConnMaxLifetime(cfg.DBConnMaxLifetime)
	db.SetConnMaxIdleTime(cfg.DBConnMaxIdleTime)

	pingDeadline := 5 * time.Second
	if cfg.ReadTimeout > 0 && cfg.ReadTimeout < pingDeadline {
		pingDeadline = cfg.ReadTimeout
	}

	if err := pingWithin(db, pingDeadline); err != nil {
		_ = db.Close()
		return nil, err
	}

	return db, nil
}

func pingWithin(db *sql.DB, timeout time.Duration) error {
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	errCh := make(chan error, 1)
	go func() {
		errCh <- db.Ping()
	}()

	select {
	case err := <-errCh:
		if err != nil {
			return fmt.Errorf("db ping: %w", err)
		}
		return nil
	case <-timer.C:
		return fmt.Errorf("db ping timeout after %s", timeout)
	}
}
