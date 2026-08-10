package integration

import (
	"testing"

	shared_db "github.com/AppeiYA/consultation-platform/internal/shared/db"
	"github.com/AppeiYA/consultation-platform/internal/shared/testhelpers"
)

func testPostgres(t *testing.T) *shared_db.Repository {
	return testhelpers.TestPostgres(t, "users")
}