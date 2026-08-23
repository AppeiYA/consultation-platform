package identity

import (
	"time"

	"github.com/AppeiYA/consultation-platform/internal/identity/ports/inbound"
	"github.com/AppeiYA/consultation-platform/internal/identity/ports/outbound"
	"github.com/AppeiYA/consultation-platform/internal/identity/usecase"
)

type Module struct {
	RegisterUser inbound.RegisterUserInt
	LoginUser    inbound.LoginUserInt
	LogoutUser   inbound.LogoutUserInt
	ValidateSession inbound.ValidateSessionInt
	GetCurrentUser inbound.GetCurrentUserInt
	UpdateUserRole inbound.UpdateUserRoleInt
}

func NewModule(
	userRepo outbound.UserRepository,
	sessionStore outbound.SessionStore,
	passwordHasher outbound.PasswordHasher,
	idGenerator outbound.IdentifierGenerator,
	clock outbound.Clock,
	sessionTokenHasher outbound.SessionTokenHasher,
	sessionTokenGenerator outbound.SessionTokenGenerator,
	sessionTTL time.Duration,
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
			sessionTTL,
		),
		LogoutUser: usecase.NewLogoutUser(
			sessionStore,
			sessionTokenHasher,
		),
		ValidateSession: usecase.NewValidateSession(
			sessionStore,
			sessionTokenHasher,
			clock,
		),
		GetCurrentUser: usecase.NewGetCurrentUser(
			userRepo,
		),
		UpdateUserRole: usecase.NewUpdateUserRoleUsecase(
			userRepo,
		),
	}
}
