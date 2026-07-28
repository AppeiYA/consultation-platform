package config

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func SetupConfig() *Config {
	app := appConfig()
	return &Config{
		App:  app,
		Http: httpConfig(),
		Session: sessionConfig(app.Env),
		Logger: loggerConfig(),
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