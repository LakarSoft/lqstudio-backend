package database

import (
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
)

// RunMigrations executes all pending database migrations from the filesystem
// migrationsDir is the path to the directory containing SQL migration files
func RunMigrations(pool *sql.DB, migrationsDir string) error {
	// Set the Goose dialect to PostgreSQL
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set goose dialect: %w", err)
	}

	// Run migrations from the filesystem
	if err := goose.Up(pool, migrationsDir); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

// GetMigrationStatus returns the current migration version
func GetMigrationStatus(pool *sql.DB) (int64, error) {
	if err := goose.SetDialect("postgres"); err != nil {
		return 0, fmt.Errorf("failed to set goose dialect: %w", err)
	}

	version, err := goose.GetDBVersion(pool)
	if err != nil {
		return 0, fmt.Errorf("failed to get migration version: %w", err)
	}

	return version, nil
}
