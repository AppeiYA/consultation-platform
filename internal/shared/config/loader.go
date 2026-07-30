package config

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func SetupConfig() *Config {
	app := appConfig()
	return &Config{
		App:      app,
		Http:     httpConfig(),
		Session:  sessionConfig(app.Env),
		Logger:   loggerConfig(),
		Database: databaseConfig(),
		Redis: redisConfig(),
	}
}

func appConfig() AppConfig {
	return AppConfig{
		Name: getEnv("APP_NAME", "CONSULTANT_SYSTEM"),
		Env:  getEnv("APP_ENV", "development"),
	}
}

func httpConfig() HttpConfig {
	return HttpConfig{
		Host: getEnv("HTTP_HOST", "localhost"),
		Port: getEnvInt("HTTP_PORT", 8080),
	}
}

func sessionConfig(app_env string) SessionConfig {
	var secure, httpOnly bool
	if app_env == "production" {
		secure, httpOnly = true, true
	} else {
		secure, httpOnly = false, true
	}
	return SessionConfig{
		CookieName: getEnv("SESSION_COOKIE_NAME", "session_id"),
		TTL:        getEnvDuration("SESSION_TTL", 24*time.Hour),
		Secure:     secure,
		HTTPOnly:   httpOnly,
		SameSite:   fiber.CookieSameSiteLaxMode,
	}
}

func loggerConfig() LoggerConfig {
	return LoggerConfig{
		LogLevel: levelFromEnv(),
	}
}

func databaseConfig() DatabaseConfig {
	return DatabaseConfig{
		Host: getEnv("DB_HOST", "localhost"),
		Port: getEnvInt("DB_PORT", 5432),
		User: getEnv("DB_USER", "peterpaul"),
		Password: getEnv("DB_PASSWORD", "Appei2004"),
		Name: getEnv("DB_NAME", "consultation_platform"),
		SSLMode: getEnv("DB_SSLMODE", "disable"),
		Schema: getEnv("DB_SCHEMA", "consultation"),
		MaxOpenConnections: getEnvInt("DB_MAX_OPEN_CONNECTIONS", 25),
		MaxIdleConnections: getEnvInt("DB_MAX_IDLE_CONNECTIONS", 5),
		MaxLifetimeMinutes: getEnvInt("DB_MAX_LIFETIME_MINUTES", 15),
	}
}

func redisConfig() RedisConfig {
	return RedisConfig{
		Address: getEnv("REDIS_ADDRESS", "localhost:6379"),
		Username: getEnv("REDIS_USERNAME", ""),
		Password: getEnv("REDIS_PASSWORD", ""),
		DB: getEnvInt("REDIS_DB", 0),
	}
}