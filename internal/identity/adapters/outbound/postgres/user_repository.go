package postgres

import (
	"github.com/AppeiYA/consultation-platform/internal/identity/ports/outbound"
	system "github.com/AppeiYA/consultation-platform/internal/shared/adapters/outbound/clock"
	"github.com/AppeiYA/consultation-platform/internal/shared/db"
)

type UserRepository struct {
	repository db.Repository
	clock *system.SystemClock
}

func NewUserRepository(repository db.Repository, clock *system.SystemClock) *UserRepository {
	return &UserRepository{
		repository: repository,
		clock: clock,
	}
}

var _ outbound.UserRepository = (*UserRepository)(nil)