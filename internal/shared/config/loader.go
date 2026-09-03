package config

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func SetupConfig() *Config {
	loadDotEnv()

	app := appConfig()
	return &Config{
		App:      app,
		Http:     httpConfig(),
		Session:  sessionConfig(app.Env),
		Logger:   loggerConfig(),
		Database: databaseConfig(),
		Redis:    redisConfig(),
		AI:       aiConfig(),
	}
}

func aiConfig() AIConfig {
	return AIConfig{
		GeminiAPIKey: getEnv("GEMINI_API_KEY", ""),
		GeminiModel:  getEnv("GEMINI_MODEL", "gemini-2.5-flash"),
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

func sessionConfig(appEnv string) SessionConfig {
    if appEnv == "production" {
        return SessionConfig{
            CookieName: getEnv("SESSION_COOKIE_NAME", "session_id"),
            TTL:        getEnvDuration("SESSION_TTL", 24*time.Hour),
            Secure:     true,
            HTTPOnly:   true,
            SameSite:   string(fiber.CookieSameSiteLaxMode),
        }
    }

    return SessionConfig{
        CookieName: getEnv("SESSION_COOKIE_NAME", "session_id"),
        TTL:        getEnvDuration("SESSION_TTL", 24*time.Hour),
        Secure:     false,
        HTTPOnly:   true,
        SameSite:   string(fiber.CookieSameSiteLaxMode),
    }
}

func loggerConfig() LoggerConfig {
	return LoggerConfig{
		LogLevel: levelFromEnv(),
	}
}

func databaseConfig() DatabaseConfig {
	return DatabaseConfig{
		Host:               requireEnv("DB_HOST"),
		Port:               requireEnvInt("DB_PORT"),
		User:               requireEnv("DB_USER"),
		Password:           requireEnv("DB_PASSWORD"),
		Name:               requireEnv("DB_NAME"),
		SSLMode:            getEnv("DB_SSLMODE", "disable"),
		Schema:             getEnv("DB_SCHEMA", "consultation"),
		MaxOpenConnections: getEnvInt("DB_MAX_OPEN_CONNECTIONS", 25),
		MaxIdleConnections: getEnvInt("DB_MAX_IDLE_CONNECTIONS", 5),
		MaxLifetimeMinutes: getEnvInt("DB_MAX_LIFETIME_MINUTES", 15),
	}
}

func redisConfig() RedisConfig {
	return RedisConfig{
		Address:  getEnv("REDIS_ADDRESS", "localhost:6379"),
		Username: getEnv("REDIS_USERNAME", ""),
		Password: getEnv("REDIS_PASSWORD", ""),
		DB:       getEnvInt("REDIS_DB", 0),
	}
}