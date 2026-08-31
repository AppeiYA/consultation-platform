package integration

import (
	"context"
	"testing"
	"time"

	"github.com/AppeiYA/consultation-platform/internal/consultant"
	consultant_http "github.com/AppeiYA/consultation-platform/internal/consultant/adapters/inbound/http"
	consultantIdentityAdapter "github.com/AppeiYA/consultation-platform/internal/consultant/adapters/outbound/external/identity"
	consultantRepo "github.com/AppeiYA/consultation-platform/internal/consultant/adapters/outbound/postgres/consultant"
	consultantAvailabilityRepo "github.com/AppeiYA/consultation-platform/internal/consultant/adapters/outbound/postgres/consultant_availability"
	consultantExpertiseRepo "github.com/AppeiYA/consultation-platform/internal/consultant/adapters/outbound/postgres/consultant_expertise"
	consultantVerificationRepo "github.com/AppeiYA/consultation-platform/internal/consultant/adapters/outbound/postgres/consultant_verification"
	professionRepo "github.com/AppeiYA/consultation-platform/internal/consultant/adapters/outbound/postgres/profession"
	consultant_verification "github.com/AppeiYA/consultation-platform/internal/consultant/adapters/outbound/verification"
	"github.com/AppeiYA/consultation-platform/internal/consultant/ports/outbound"
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

type ConsultantAppOption func(*consultantAppConfig)

type consultantAppConfig struct {
	verificationService outbound.VerificationService
}

func WithVerificationService(service outbound.VerificationService) ConsultantAppOption {
	return func(c *consultantAppConfig) {
		c.verificationService = service
	}
}

func setUpConsultantApp(t *testing.T, opts ...ConsultantAppOption) *TestHarness {
	t.Helper()
	cfg := &consultantAppConfig{}
	for _, opt := range opts {
		opt(cfg)
	}

	appConfig := config.SetupTestConfig()
	db := testhelpers.TestPostgres(t, "users", "consultants", "consultant_verifications", "consultant_availabilities", "consultant_expertises")
	redisClient := testhelpers.TestRedis(t)

	clock := system.NewSystemClock()
	idGenerator := uuid.NewGenerator()

	// Seed professions if not present
	ctx := context.Background()
	_, _ = db.Executor(ctx).ExecContext(ctx, `
		INSERT INTO professions (id, name, created_at) VALUES
			('prof_9ee432d7-b672-40ae-b03f-c1f1fb696621', 'SOFTWARE_ENGINEER', CURRENT_TIMESTAMP),
			('prof_12d965f5-e1f5-49aa-ac57-856772d236ce', 'LAWYER', CURRENT_TIMESTAMP),
			('prof_940f840d-617a-4ead-8a8c-873851762bc7', 'DOCTOR', CURRENT_TIMESTAMP),
			('prof_9a1e78f7-a053-4b4b-8802-fcac77e530ec', 'ACCOUNTANT', CURRENT_TIMESTAMP),
			('prof_d95d5c58-d5be-4bca-bf87-b84b3f2a2681', 'THERAPIST', CURRENT_TIMESTAMP),
			('prof_aef03e88-a0ee-49c9-b455-4b2210412b52', 'CLERGY', CURRENT_TIMESTAMP)
		ON CONFLICT (id) DO NOTHING;
	`)

	// Consultant repos
	consultantRepository := consultantRepo.NewConsultantRepository(*db, clock)
	professionRepository := professionRepo.NewProfessionRepository(*db, clock)
	verificationRepo := consultantVerificationRepo.NewVerificationRepository(*db, clock)
	availabilityRepo := consultantAvailabilityRepo.NewAvailabilityRepository(*db, clock)
	expertiseRepo := consultantExpertiseRepo.NewConsultantExpertiseRepository(db)

	var verificationService outbound.VerificationService = &consultant_verification.UnavailableVerificationService{}
	if cfg.verificationService != nil {
		verificationService = cfg.verificationService
	}

	// Identity dependencies
	userRepo := identity_postgres.NewUserRepository(*db, clock)
	sessionStore := identity_redis.NewSessionStore(redisClient, clock)
	passwordHasher := identity_bcrypt.NewHasher()
	sessionTokenHasher := identity_sha256.NewHasher()
	sessionTokenGenerator := identity_random.NewGenerator()
	cookieManager := shared_http.NewCookieManager(appConfig.Session)

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
	roleAssigner := consultantIdentityAdapter.NewRoleAssigner(identityModule)
	consultantModule := consultant.NewModule(
		consultantRepository,
		verificationService,
		verificationRepo,
		availabilityRepo,
		professionRepository,
		expertiseRepo,
		roleAssigner,
		idGenerator,
		clock,
	)

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
