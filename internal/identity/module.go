package identity

import (
	"github.com/AppeiYA/consultation-platform/internal/identity/ports/inbound"
	"github.com/AppeiYA/consultation-platform/internal/identity/ports/outbound"
	"github.com/AppeiYA/consultation-platform/internal/identity/usecase"
)

type Module struct {
	RegisterUser inbound.RegisterUserInt
	LoginUser    inbound.LoginUserInt
	LogoutUser   inbound.LogoutUserInt
	GetCurrentUser inbound.GetCurrentUserInt
}

func NewModule(
	userRepo outbound.UserRepository,
	sessionStore outbound.SessionStore,
	passwordHasher outbound.PasswordHasher,
	idGenerator outbound.IdentifierGenerator,
	clock outbound.Clock,
	sessionTokenHasher outbound.SessionTokenHasher,
	sessionTokenGenerator outbound.SessionTokenGenerator,
) *Module {
	return &Module{
		RegisterUser: usecase.NewRegisterUser(
			userRepo,
			passwordHasher,
			idGenerator,
			clock,
		),
		LoginUser: usecase.NewLoginUser(
			userRepo,
			sessionStore,
			passwordHasher,
			sessionTokenHasher,
			sessionTokenGenerator,
			idGenerator,
			clock,
		),
		LogoutUser: usecase.NewLogoutUser(
			sessionStore,
			sessionTokenHasher,
		),
		GetCurrentUser: usecase.NewGetCurrentUser(
			sessionStore,
			sessionTokenHasher,
			userRepo,
			clock,
		),
	}
}
