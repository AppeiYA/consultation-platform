package http

import (
	"github.com/AppeiYA/consultation-platform/internal/identity/adapters/inbound/http/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetUpRouter(router fiber.Router, h *IdentityHandler, m *middleware.AuthenticationMiddleware) {
	identity := router.Group("/identity")

	identity.Post("/register", h.Register)
	identity.Post("/login", h.Login)

	protected := identity.Use(m.Authenticate)

	protected.Get("/me", h.GetMe)
}