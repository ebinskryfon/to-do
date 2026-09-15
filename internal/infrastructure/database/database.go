package database

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"todo/pkg/config"
)

// NewGORM initializes a GORM PostgreSQL database connection with connection pooling
// and a bounded connection ping, following the Ollie Cake House architecture pattern.
func NewGORM(cfg *config.DatabaseConfig, log zerolog.Logger) (*gorm.DB, error) {
	dsn := cfg.DSN()
	// Ensure the dial itself is bounded; pgx stdlib honors connect_timeout (seconds).
	if !strings.Contains(dsn, "connect_timeout=") {
		if strings.Contains(dsn, "?") {
			dsn += "&connect_timeout=20"
		} else {
			dsn += "?connect_timeout=20"
		}
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB from gorm: %w", err)
	}

	maxOpen := cfg.MaxOpenConns
	if maxOpen <= 0 {
		maxOpen = 10
	}
	maxIdle := cfg.MaxIdleConns
	if maxIdle <= 0 {
		maxIdle = 2
	}
	maxIdleTime := cfg.ConnMaxIdleTime
	if maxIdleTime <= 0 {
		maxIdleTime = 30 * time.Second
	}
	maxLifetime := cfg.ConnMaxLifetime
	if maxLifetime <= 0 {
		maxLifetime = 5 * time.Minute
	}

	// Database Connection Pool settings
	sqlDB.SetMaxOpenConns(maxOpen)
	sqlDB.SetMaxIdleConns(maxIdle)
	sqlDB.SetConnMaxIdleTime(maxIdleTime)
	sqlDB.SetConnMaxLifetime(maxLifetime)

	log.Info().
		Int("max_open_conns", maxOpen).
		Int("max_idle_conns", maxIdle).
		Dur("conn_max_idle_time", maxIdleTime).
		Dur("conn_max_lifetime", maxLifetime).
		Msg("database connection pool configured")

	pingCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(pingCtx); err != nil {
		return nil, fmt.Errorf("database unreachable within 20s: %w", err)
	}

	return db, nil
}

// NewPostgresConnection provides a convenient wrapper taking the root *config.Config.
func NewPostgresConnection(cfg *config.Config, log zerolog.Logger) (*gorm.DB, error) {
	return NewGORM(&cfg.Database, log)
}
