package database

import (
	"context"
	"database/sql"
	"fmt"
	"lqstudio-backend/internal/config"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib" // PostgreSQL driver for database/sql
)

// Connection represents a database connection pool
type Connection struct {
	Pool   *pgxpool.Pool
	config *config.DatabaseConfig
}

// New creates a new database connection pool
func New(cfg *config.DatabaseConfig) (*Connection, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var pgxCfg *pgxpool.Config
	var err error

	// Use DATABASE_URL if available (Neon-compatible)
	if cfg.DatabaseURL != "" {
		pgxCfg, err = pgxpool.ParseConfig(cfg.DatabaseURL)
		if err != nil {
			return nil, fmt.Errorf("failed to parse DATABASE_URL: %w", err)
		}
	} else {
		// Build connection string from individual fields
		connStr := cfg.ConnectionString()
		pgxCfg, err = pgxpool.ParseConfig(connStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse connection config: %w", err)
		}
	}

	// Ensure schema and timezone are set in runtime params
	if pgxCfg.ConnConfig.RuntimeParams == nil {
		pgxCfg.ConnConfig.RuntimeParams = make(map[string]string)
	}
	if _, exists := pgxCfg.ConnConfig.RuntimeParams["search_path"]; !exists && cfg.Schema != "" {
		pgxCfg.ConnConfig.RuntimeParams["search_path"] = cfg.Schema
	}
	if _, exists := pgxCfg.ConnConfig.RuntimeParams["timezone"]; !exists && cfg.Timezone != "" {
		pgxCfg.ConnConfig.RuntimeParams["timezone"] = cfg.Timezone
	}

	// Apply connection pool settings from config
	pgxCfg.MaxConns = cfg.MaxConns
	pgxCfg.MinConns = cfg.MinConns
	pgxCfg.MaxConnLifetime = cfg.MaxConnLifetime
	pgxCfg.MaxConnIdleTime = cfg.MaxConnIdleTime

	pool, err := pgxpool.NewWithConfig(ctx, pgxCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Test connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Connection{
		Pool:   pool,
		config: cfg,
	}, nil

}

// Close closes the database connection pool
func (c *Connection) Close() {
	c.Pool.Close()
}

// Health checks the health of the database connection
func (c *Connection) Health(ctx context.Context) error {
	return c.Pool.Ping(ctx)
}

// GetStdDB returns a standard database/sql connection for migration tools
// Note: This creates a new connection - caller is responsible for closing it
func (c *Connection) GetStdDB() (*sql.DB, error) {
	// Use DATABASE_URL if available, otherwise build connection string
	var connStr string
	if c.config.DatabaseURL != "" {
		connStr = c.config.DatabaseURL
	} else {
		connStr = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s search_path=%s",
			c.config.Host,
			c.config.Port,
			c.config.Username,
			c.config.Password,
			c.config.Database,
			c.config.SSLMode,
			c.config.Schema,
		)
	}

	// Open connection using pgx driver
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database/sql connection: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}
