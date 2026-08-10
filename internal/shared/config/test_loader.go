package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// TestConfigOption defines a functional option for customizing test configuration.
type TestConfigOption func(*Config)

// SetupTestConfig constructs a dedicated configuration for integration tests.
// It defaults to isolated resources (e.g., consultation_platform_test DB, Redis DB index 15)
// and enforces strict safety checks to prevent tests from executing against production resources.
func SetupTestConfig(opts ...TestConfigOption) *Config {
	loadDotEnv()

	cfg := &Config{
		App: AppConfig{
			Name: getEnv("TEST_APP_NAME", "CONSULTANT_SYSTEM_TEST"),
			Env:  "test",
		},
		Http: HttpConfig{
			Host: getEnv("TEST_HTTP_HOST", "localhost"),
			Port: getEnvInt("TEST_HTTP_PORT", 0),
		},
		Session: SessionConfig{
			CookieName: getEnv("TEST_SESSION_COOKIE_NAME", "test_session_id"),
			TTL:        getEnvDuration("TEST_SESSION_TTL", 24*time.Hour),
			Secure:     false,
			HTTPOnly:   true,
			SameSite:   fiber.CookieSameSiteLaxMode,
		},
		Logger: LoggerConfig{
			LogLevel: levelFromEnv(),
		},
		Database: DatabaseConfig{
			Host:               getEnv("TEST_DB_HOST", getEnv("DB_HOST", "localhost")),
			Port:               getEnvInt("TEST_DB_PORT", getEnvInt("DB_PORT", 5432)),
			User:               getEnv("TEST_DB_USER", getEnv("DB_USER", "peterpaul")),
			Password:           getEnv("TEST_DB_PASSWORD", getEnv("DB_PASSWORD", "")),
			Name:               getEnv("TEST_DB_NAME", "consultation_platform_test"),
			SSLMode:            getEnv("TEST_DB_SSLMODE", getEnv("DB_SSLMODE", "disable")),
			Schema:             getEnv("TEST_DB_SCHEMA", getEnv("DB_SCHEMA", "consultation")),
			MaxOpenConnections: getEnvInt("TEST_DB_MAX_OPEN_CONNECTIONS", 10),
			MaxIdleConnections: getEnvInt("TEST_DB_MAX_IDLE_CONNECTIONS", 5),
			MaxLifetimeMinutes: getEnvInt("TEST_DB_MAX_LIFETIME_MINUTES", 5),
		},
		Redis: RedisConfig{
			Address:  getEnv("TEST_REDIS_ADDRESS", getEnv("REDIS_ADDRESS", "localhost:6379")),
			Username: getEnv("TEST_REDIS_USERNAME", ""),
			Password: getEnv("TEST_REDIS_PASSWORD", ""),
			DB:       getEnvInt("TEST_REDIS_DB", 15),
		},
	}

	for _, opt := range opts {
		opt(cfg)
	}

	if err := validateTestSafety(cfg); err != nil {
		panic(fmt.Sprintf("CRITICAL TEST SAFETY VIOLATION: %v", err))
	}

	return cfg
}

// validateTestSafety enforces that test configuration explicitly points to isolated test resources.
// It fails fast if the environment, database name, Redis DB index, or session cookie name
// do not indicate a dedicated test environment.
func validateTestSafety(cfg *Config) error {
	env := strings.ToLower(cfg.App.Env)
	if env != "test" && env != "testing" {
		return fmt.Errorf("refusing to run integration tests with APP_ENV=%q (must be 'test' or 'testing')", cfg.App.Env)
	}

	if !strings.Contains(strings.ToLower(cfg.Database.Name), "test") {
		return fmt.Errorf("refusing to run integration tests against database %q (database name must contain 'test', e.g. '_test')", cfg.Database.Name)
	}

	prodDB := getEnv("DB_NAME", "consultation_platform")
	if cfg.Database.Name == prodDB && !strings.Contains(prodDB, "test") {
		return fmt.Errorf("refusing to run integration tests against production/dev database %q", cfg.Database.Name)
	}

	if cfg.Redis.DB == 0 {
		return fmt.Errorf("refusing to run integration tests against Redis DB 0 (must use dedicated non-zero test DB index, e.g. 15)")
	}

	if !strings.Contains(strings.ToLower(cfg.Session.CookieName), "test") {
		return fmt.Errorf("refusing to run integration tests with session cookie %q (cookie name must contain 'test')", cfg.Session.CookieName)
	}

	return nil
}

// WithTestDatabaseName allows overriding the test database name.
func WithTestDatabaseName(name string) TestConfigOption {
	return func(c *Config) {
		c.Database.Name = name
	}
}

// WithTestRedisDB allows overriding the test Redis database index.
func WithTestRedisDB(db int) TestConfigOption {
	return func(c *Config) {
		c.Redis.DB = db
	}
}
