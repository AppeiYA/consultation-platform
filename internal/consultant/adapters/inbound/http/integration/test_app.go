package integration

import (
	"testing"
	"time"

	"github.com/AppeiYA/consultation-platform/internal/consultant"
	consultant_http "github.com/AppeiYA/consultation-platform/internal/consultant/adapters/inbound/http"
	consultant_postgres "github.com/AppeiYA/consultation-platform/internal/consultant/adapters/outbound/postgres"
	"github.com/AppeiYA/consultation-platform/internal/identity"
	identity_http "github.com/AppeiYA/consultation-platform/internal/identity/adapters/inbound/http"
	identity_middleware "github.com/AppeiYA/consultation-platform/internal/identity/adapters/inbound/http/middleware"
	identity_bcrypt "github.com/AppeiYA/consultation-platform/internal/identity/adapters/outbound/bcrypt"
	identity_postgres "github.com/AppeiYA/consultation-platform/internal/identity/adapters/outbound/postgres"
	identity_random "github.com/AppeiYA/consultation-platform/internal/identity/adapters/outbound/random"
	identity_redis "github.com/AppeiYA/consultation-platform/internal/identity/adapters/outbound/redis"
	identity_sha256 "github.com/AppeiYA/consultation-platform/internal/identity/adapters/outbound/sha256"
	shared_http "github.com/AppeiYA/consultation-platform/internal/shared/adapters/http"
	"github.com/AppeiYA/consultation-platform/internal/shared/adapters/id/uuid"
	system "github.com/AppeiYA/consultation-platform/internal/shared/adapters/outbound/clock"
	"github.com/AppeiYA/consultation-platform/internal/shared/config"
	"github.com/AppeiYA/consultation-platform/internal/shared/testhelpers"
	"github.com/gofiber/fiber/v2"
)

type TestHarness struct {
	App *fiber.App
}

func setUpConsultantApp(t *testing.T) *TestHarness {
	t.Helper()
	cfg := config.SetupTestConfig()

	// Infrastructure using shared testhelpers
	db := testhelpers.TestPostgres(t, "users", "consultants")
	rdb := testhelpers.TestRedis(t)

	// Shared services
	cookieManager := shared_http.NewCookieManager(cfg.Session)
	idGenerator := uuid.NewGenerator()
	clock := system.NewSystemClock()
	passwordHasher := identity_bcrypt.NewHasher()
	sessionTokenHasher := identity_sha256.NewHasher()
	sessionTokenGenerator := identity_random.NewGenerator()

	// Repositories
	userRepo := identity_postgres.NewUserRepository(*db, clock)
	sessionStore := identity_redis.NewSessionStore(rdb, clock)
	consultantRepo := consultant_postgres.NewConsultantRepository(*db, clock)

	// Modules
	identityModule := identity.NewModule(
		userRepo,
		sessionStore,
		passwordHasher,
		idGenerator,
		clock,
		sessionTokenHasher,
		sessionTokenGenerator,
		time.Hour,
	)
	consultantModule := consultant.NewModule(consultantRepo, idGenerator, clock)

	// Handlers & Middleware
	identityHandler := identity_http.NewIdentityHandler(identityModule, cookieManager)
	identityMiddleware := identity_middleware.NewAuthenticationMiddleware(identityModule, cookieManager)
	consultantHandler := consultant_http.NewConsultantHandler(*consultantModule)

	// Fiber app
	app := fiber.New()
	testGroup := app.Group("/test/v1")

	// Routes
	identity_http.SetUpRouter(testGroup, identityHandler, identityMiddleware)
	consultant_http.RegisterConsultantRoutes(testGroup, consultantHandler, identityMiddleware)

	return &TestHarness{
		App: app,
	}
}
