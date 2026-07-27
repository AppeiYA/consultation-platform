package http

import "github.com/gofiber/fiber/v2"

func SetUpRouter(router fiber.Router, h *IdentityHandler) {
	identity := router.Group("/identity")

	identity.Post("/register", h.Register)
}