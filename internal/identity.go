package app

import (
	"time"

	identity_module "github.com/AppeiYA/consultation-platform/internal/identity"
	identity_http "github.com/AppeiYA/consultation-platform/internal/identity/adapters/inbound/http"
	identity_auth_middleware "github.com/AppeiYA/consultation-platform/internal/identity/adapters/inbound/http/middleware"
	"github.com/AppeiYA/consultation-platform/internal/identity/adapters/outbound/bcrypt"
	identity_memory "github.com/AppeiYA/consultation-platform/internal/identity/adapters/outbound/memory"
	"github.com/AppeiYA/consultation-platform/internal/identity/adapters/outbound/random"
	"github.com/AppeiYA/consultation-platform/internal/identity/adapters/outbound/sha256"
	shared_http "github.com/AppeiYA/consultation-platform/internal/shared/adapters/http"
	"github.com/AppeiYA/consultation-platform/internal/shared/adapters/id/uuid"
	system "github.com/AppeiYA/consultation-platform/internal/shared/adapters/outbound/clock"
)

func (a *App) registerIdentity(
	clock *system.SystemClock,
	idGenerator *uuid.Generator,
	cookieManager shared_http.CookieManagerInt,
	sessionTTL time.Duration,
) {
	// Identity
	userRepo := identity_memory.NewUserRepository()
	sessionStore := identity_memory.NewSessionStore()

	passwordHasher := bcrypt.NewHasher()
	sessionTokenHasher := sha256.NewHasher()
	sessionTokenGenerator := random.NewGenerator()

	identityModule := identity_module.NewModule(
		userRepo,
		sessionStore,
		passwordHasher,
		idGenerator,
		clock,
		sessionTokenHasher,
		sessionTokenGenerator,
		sessionTTL,
	)

	identityHandler := identity_http.NewIdentityHandler(
		identityModule, 
		cookieManager,
	)

	identityAuthMiddleware := identity_auth_middleware.NewAuthenticationMiddleware(
		identityModule, 
		cookieManager,
	)

	a.identityModule = identityModule
	a.identityHandler = identityHandler
	a.identityAuthMiddleware = identityAuthMiddleware
}