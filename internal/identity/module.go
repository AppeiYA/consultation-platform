package identity

import (
	"github.com/AppeiYA/consultation-platform/internal/identity/ports/inbound"
	"github.com/AppeiYA/consultation-platform/internal/identity/ports/outbound"
	"github.com/AppeiYA/consultation-platform/internal/identity/usecase"
)

type Module struct {
	RegisterUser inbound.RegisterUserInt
}

func NewModule(
	userRepo outbound.UserRepository,
	sessionStore outbound.SessionStore,
	passwordHasher outbound.PasswordHasher,
	idGenerator outbound.IdentifierGenerator,
	clock outbound.Clock,
	sessionTokenHasher outbound.SessionTokenHasher,
) *Module {
	return &Module{
		RegisterUser: usecase.NewRegisterUser(
			userRepo,
			passwordHasher,
			idGenerator,
			clock,
		),
	}
}
