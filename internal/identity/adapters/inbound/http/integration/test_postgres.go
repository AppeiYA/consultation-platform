package integration

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/AppeiYA/consultation-platform/internal/shared/config"
	shared_db "github.com/AppeiYA/consultation-platform/internal/shared/db"
)

var tables = []string{"users"}

func testPostgres(t *testing.T) *shared_db.Repository {
	t.Helper()

	cfg := config.SetupConfig()

	db, err := shared_db.Connect(cfg.Database)
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}

	cleanupPostgres(t, db)

	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Errorf("failed to close database: %v", err)
		}
	})

	repo := shared_db.NewRepository(db)

	return &repo
}

func cleanupPostgres(t *testing.T, db *shared_db.DB) {
	t.Helper()

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