package testhelpers

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/AppeiYA/consultation-platform/internal/shared/config"
	shared_db "github.com/AppeiYA/consultation-platform/internal/shared/db"
	"github.com/golang-migrate/migrate/v4"
	migratepostgres "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func TestPostgres(t *testing.T, tables ...string) *shared_db.Repository {
	t.Helper()

	cfg := config.SetupTestConfig()

	db, err := shared_db.Connect(cfg.Database)
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}

	runTestMigrations(t, db, cfg.Database)
	CleanupPostgres(t, db, tables...)

	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Errorf("failed to close database: %v", err)
		}
	})

	repo := shared_db.NewRepository(db)

	return &repo
}

func runTestMigrations(t *testing.T, db *shared_db.DB, cfg config.DatabaseConfig) {
	t.Helper()

	migPath := findMigrationsDir()
	if migPath == "" {
		t.Fatalf("failed to locate migrations directory")
	}

	ctx := context.Background()
	executor := db.Executor(ctx)

	if cfg.Schema != "" && cfg.Schema != "public" {
		_, err := executor.ExecContext(ctx, fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s;", cfg.Schema))
		if err != nil {
			t.Fatalf("failed to create schema %s: %v", cfg.Schema, err)
		}
	}

	driver, err := migratepostgres.WithInstance(db.Conn().DB, &migratepostgres.Config{
		DatabaseName: cfg.Name,
		SchemaName:   cfg.Schema,
	})
	if err != nil {
		t.Fatalf("failed to create migration driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migPath),
		cfg.Name,
		driver,
	)
	if err != nil {
		t.Fatalf("failed to create migrate instance: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		if strings.Contains(err.Error(), "Dirty database") {
			_ = m.Force(1)
			if err := m.Up(); err != nil && err != migrate.ErrNoChange {
				t.Fatalf("failed to run migrations on test database after clearing dirty state: %v", err)
			}
		} else {
			t.Fatalf("failed to run migrations on test database: %v", err)
		}
	}
}

func findMigrationsDir() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}

	for {
		migPath := filepath.Join(dir, "migrations")
		if info, err := os.Stat(migPath); err == nil && info.IsDir() {
			return migPath
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return ""
}

func CleanupPostgres(t *testing.T, db *shared_db.DB, tables ...string) {
	t.Helper()

	if len(tables) == 0 {
		return
	}

	ctx := context.Background()
	executor := db.Executor(ctx)

	_, err := executor.ExecContext(ctx, fmt.Sprintf(`
		TRUNCATE TABLE
			%s
		RESTART IDENTITY CASCADE
	`, strings.Join(tables, ", ")))
	if err != nil {
		t.Fatalf("failed to clean database: %v", err)
	}
}
