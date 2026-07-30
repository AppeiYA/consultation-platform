package db

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/AppeiYA/consultation-platform/internal/shared/config"
)

func Connect(cfg config.DatabaseConfig) (*DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s search_path=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.Name,
		cfg.SSLMode,
		cfg.Schema,
	)

	conn, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}

	conn.SetMaxOpenConns(cfg.MaxOpenConnections)
	conn.SetMaxIdleConns(cfg.MaxIdleConnections)
	conn.SetConnMaxLifetime(time.Duration(cfg.MaxLifetimeMinutes) * time.Minute)

	return New(conn), nil
}