package app

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AppeiYA/consultation-platform/internal/identity"
	identity_http "github.com/AppeiYA/consultation-platform/internal/identity/adapters/inbound/http"
	identity_auth_middleware "github.com/AppeiYA/consultation-platform/internal/identity/adapters/inbound/http/middleware"
	shared_http "github.com/AppeiYA/consultation-platform/internal/shared/adapters/http"
	"github.com/AppeiYA/consultation-platform/internal/shared/adapters/id/uuid"
	system "github.com/AppeiYA/consultation-platform/internal/shared/adapters/outbound/clock"
	"github.com/AppeiYA/consultation-platform/internal/shared/config"
	"github.com/AppeiYA/consultation-platform/internal/shared/db"
	"github.com/AppeiYA/consultation-platform/internal/shared/logger"
	"github.com/AppeiYA/consultation-platform/internal/shared/redis"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type App struct {
	config *config.Config
	fiber *fiber.App

	identityModule *identity.Module
	identityHandler *identity_http.IdentityHandler
	identityAuthMiddleware *identity_auth_middleware.AuthenticationMiddleware
}

func New() (*App, error) {
	cfg := config.SetupConfig()
	database, err := db.Connect(cfg.Database)
	if err != nil {
		return nil, err
	}

	redisPK, err := redis.Connect(cfg.Redis)
	if err != nil {
		return nil, err
	}

	logger.Init(cfg)

	a := &App{
		config: cfg,
		fiber:  fiber.New(fiber.Config{
			JSONEncoder: func(v interface{}) ([]byte, error) {
				buf := &bytes.Buffer{}
				encoder := json.NewEncoder(buf)
				encoder.SetEscapeHTML(false)
				err := encoder.Encode(v)
				return bytes.TrimRight(buf.Bytes(), "\n"), err
			},
			EnableTrustedProxyCheck: true,
			TrustedProxies:          []string{"0.0.0.0/0"},
			BodyLimit:               5 * 1024 * 1024,
		}),
	}

	a.fiber.Use(cors.New(cors.Config{
		AllowOriginsFunc: func(origin string) bool {
			return true
		},
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowCredentials: true,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, session_id",
	}))
	a.fiber.Use(recover.New())

	a.fiber.Options("/*", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNoContent)
	})

	// shared
	cookieManager := shared_http.NewCookieManager(cfg.Session)
	idGenerator := uuid.NewGenerator()
	clock := system.NewSystemClock()
	repository := db.NewRepository(database)

	a.registerIdentity(
		repository,
		redisPK,
		clock,
		idGenerator,
		cookieManager,
		cfg.Session.TTL,
	)

	SetUpRouter(a)

	return a, nil
}

func (a *App) Run() error {
	errCh := make(chan error, 1)
	go func() {
		if err := a.fiber.Listen(":" + fmt.Sprintf("%d", a.config.Http.Port)); err != nil {
			errCh <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-errCh:
		return fmt.Errorf("server failed to start: %w", err)
	case <-quit:
		logger.Info("shutting down gracefully...")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := a.fiber.ShutdownWithContext(ctx); err != nil {
		return fmt.Errorf("graceful shutdown failed: %w", err)
	}

	logger.Info("shutdown complete")
	return nil
}



