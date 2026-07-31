package integration

import (
	"testing"
	"time"

	"github.com/AppeiYA/consultation-platform/internal/identity"
	identity_http "github.com/AppeiYA/consultation-platform/internal/identity/adapters/inbound/http"
	identity_middleware "github.com/AppeiYA/consultation-platform/internal/identity/adapters/inbound/http/middleware"
	"github.com/AppeiYA/consultation-platform/internal/identity/adapters/outbound/bcrypt"
	"github.com/AppeiYA/consultation-platform/internal/identity/adapters/outbound/postgres"
	"github.com/AppeiYA/consultation-platform/internal/identity/adapters/outbound/random"
	"github.com/AppeiYA/consultation-platform/internal/identity/adapters/outbound/redis"
	"github.com/AppeiYA/consultation-platform/internal/identity/adapters/outbound/sha256"
	shared_http "github.com/AppeiYA/consultation-platform/internal/shared/adapters/http"
	"github.com/AppeiYA/consultation-platform/internal/shared/adapters/id/uuid"
	system "github.com/AppeiYA/consultation-platform/internal/shared/adapters/outbound/clock"
	"github.com/AppeiYA/consultation-platform/internal/shared/config"
	"github.com/gofiber/fiber/v2"
)

type TestHarness struct {
	App *fiber.App
}

func setUpIdentityApp(t *testing.T) *TestHarness {
	t.Helper()
	cfg := config.SetupTestConfig()

	// Infrastructure
	db := testPostgres(t)
	rdb := testRedis(t)

	// Shared services
	cookieManager := shared_http.NewCookieManager(cfg.Session)
	idGenerator := uuid.NewGenerator()
	clock := system.NewSystemClock()
	passwordHasher := bcrypt.NewHasher()
	sessionTokenHasher := sha256.NewHasher()
	sessionTokenGenerator := random.NewGenerator()

	// Repositories
	userRepo := postgres.NewUserRepository(*db, clock)
	sessionStore := redis.NewSessionStore(rdb, clock)

	// Module
	module := identity.NewModule(userRepo, sessionStore, passwordHasher, idGenerator, clock, sessionTokenHasher, sessionTokenGenerator, time.Hour)

	// Handler
	handler := identity_http.NewIdentityHandler(module, cookieManager)

	// Middleware
	middleware := identity_middleware.NewAuthenticationMiddleware(module, cookieManager)

	// Fiber
	app := fiber.New()

	// Routes
	setUpRoutes(app, handler, middleware)

	return &TestHarness{
		App: app,
	}
}

func setUpRoutes(
	app *fiber.App, 
	h *identity_http.IdentityHandler, 
	m *identity_middleware.AuthenticationMiddleware,
) {
	test := app.Group("/test/v1")
	identity_http.SetUpRouter(test, h, m)
}