package integration

import (
	"net/http"
	"testing"
	"time"

	"github.com/AppeiYA/consultation-platform/internal/consultationcase"
	consultationcase_http "github.com/AppeiYA/consultation-platform/internal/consultationcase/adapters/inbound/http"
	consultationCaseIdentityAdapter "github.com/AppeiYA/consultation-platform/internal/consultationcase/adapters/outbound/external/identity"
	consultationCaseRepo "github.com/AppeiYA/consultation-platform/internal/consultationcase/adapters/outbound/postgres/consultationcase"
	"github.com/AppeiYA/consultation-platform/internal/identity"
	identity_http "github.com/AppeiYA/consultation-platform/internal/identity/adapters/inbound/http"
	identity_dto "github.com/AppeiYA/consultation-platform/internal/identity/adapters/inbound/http/dto"
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
	"github.com/stretchr/testify/require"
)

type TestHarness struct {
	App *fiber.App
}

func setUpConsultationCaseApp(t *testing.T) *TestHarness {
	t.Helper()
	cfg := config.SetupTestConfig()

	// Infrastructure using shared testhelpers
	db := testhelpers.TestPostgres(t, "users", "consultation_cases")
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
	caseRepo := consultationCaseRepo.NewConsultationCaseRepository(db)

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
	clientVerifier := consultationCaseIdentityAdapter.NewClientVerifier(identityModule)
	caseModule := consultationcase.NewModule(
		caseRepo,
		clientVerifier,
		idGenerator,
		clock,
	)

	// Handlers & Middleware
	identityHandler := identity_http.NewIdentityHandler(identityModule, cookieManager)
	identityMiddleware := identity_middleware.NewAuthenticationMiddleware(identityModule, cookieManager)
	caseHandler := consultationcase_http.NewConsultationCaseHandler(caseModule)

	// Fiber app
	app := fiber.New()
	testGroup := app.Group("/test/v1")

	// Routes
	identity_http.SetUpRouter(testGroup, identityHandler, identityMiddleware)
	consultationcase_http.RegisterConsultationCaseRoutes(testGroup, caseHandler, identityMiddleware)

	return &TestHarness{
		App: app,
	}
}

func registerAndLoginUser(t *testing.T, harness *TestHarness, email string) string {
	t.Helper()

	regBody := identity_dto.RegisterRequest{
		FirstName: "John",
		LastName:  "Doe",
		Email:     email,
		Password:  "SecurePassword123!",
	}

	regResp, err := testhelpers.PerformRequest(
		harness.App,
		http.MethodPost,
		"/test/v1/identity/register",
		regBody,
	)
	require.NoError(t, err)
	regResp.Body.Close()
	require.Equal(t, http.StatusCreated, regResp.StatusCode)

	loginBody := identity_dto.LoginRequest{
		Email:    email,
		Password: "SecurePassword123!",
	}

	loginResp, err := testhelpers.PerformRequest(
		harness.App,
		http.MethodPost,
		"/test/v1/identity/login",
		loginBody,
	)
	require.NoError(t, err)
	defer loginResp.Body.Close()
	require.Equal(t, http.StatusOK, loginResp.StatusCode)

	cookieHeader := testhelpers.ExtractCookieHeader(loginResp)
	require.NotEmpty(t, cookieHeader)

	return cookieHeader
}
