package app

import (
	"time"

	identity_module "github.com/AppeiYA/consultation-platform/internal/identity"
	identity_http "github.com/AppeiYA/consultation-platform/internal/identity/adapters/inbound/http"
	identity_auth_middleware "github.com/AppeiYA/consultation-platform/internal/identity/adapters/inbound/http/middleware"
	"github.com/AppeiYA/consultation-platform/internal/identity/adapters/outbound/bcrypt"
	identity_postgres "github.com/AppeiYA/consultation-platform/internal/identity/adapters/outbound/postgres"
	"github.com/AppeiYA/consultation-platform/internal/identity/adapters/outbound/random"
	"github.com/AppeiYA/consultation-platform/internal/identity/adapters/outbound/redis"
	"github.com/AppeiYA/consultation-platform/internal/identity/adapters/outbound/sha256"
	shared_http "github.com/AppeiYA/consultation-platform/internal/shared/adapters/http"
	"github.com/AppeiYA/consultation-platform/internal/shared/adapters/id/uuid"
	system "github.com/AppeiYA/consultation-platform/internal/shared/adapters/outbound/clock"
	"github.com/AppeiYA/consultation-platform/internal/shared/db"
	shared_redis "github.com/AppeiYA/consultation-platform/internal/shared/redis"
)

func (a *App) registerIdentity(
	db db.Repository,
	redisPK *shared_redis.Redis,
	clock *system.SystemClock,
	idGenerator *uuid.Generator,
	cookieManager shared_http.CookieManagerInt,
	sessionTTL time.Duration,
) {
	// Identity
	userRepo := identity_postgres.NewUserRepository(db, clock)
	sessionStore := redis.NewSessionStore(redisPK, clock)

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