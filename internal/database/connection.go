package database

import (
	"context"
	"database/sql"
	"fmt"
	"lqstudio-backend/internal/config"
	"strconv"
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

	pgxCfg, err := pgxpool.ParseConfig("")
	if err != nil {
		return nil, fmt.Errorf("failed to parse pgx config: %w", err)
	}

	// Parse port string to integer
	port, err := strconv.Atoi(cfg.Port)
	if err != nil {
		return nil, fmt.Errorf("invalid port number: %w", err)
	}

	pgxCfg.ConnConfig.Host = cfg.Host
	pgxCfg.ConnConfig.Port = uint16(port)
	pgxCfg.ConnConfig.User = cfg.Username
	pgxCfg.ConnConfig.Password = cfg.Password
	pgxCfg.ConnConfig.Database = cfg.Database

	pgxCfg.ConnConfig.RuntimeParams = map[string]string{
		"search_path": cfg.Schema,
		"timezone":    cfg.Timezone,
	}

	// Optional but recommended pool tuning
	pgxCfg.MaxConns = 20
	pgxCfg.MinConns = 2
	pgxCfg.MaxConnLifetime = time.Hour
	pgxCfg.MaxConnIdleTime = 30 * time.Minute

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
	// Build connection string for database/sql
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable search_path=%s",
		c.config.Host,
		c.config.Port,
		c.config.Username,
		c.config.Password,
		c.config.Database,
		c.config.Schema,
	)

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
