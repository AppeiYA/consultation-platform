package http

import (
	identity_auth_middleware "github.com/AppeiYA/consultation-platform/internal/identity/adapters/inbound/http/middleware"
	"github.com/gofiber/fiber/v2"
)

func RegisterExpertMatchingRoutes(
	router fiber.Router,
	handler *ExpertMatchingHandler,
	authMiddleware *identity_auth_middleware.AuthenticationMiddleware,
) {
	group := router.Group("/expert-matching")
	group.Use(authMiddleware.Authenticate)

	group.Post("cases/:id/match", handler.StartMatching)
	group.Get("cases/:id/matches", handler.GetCaseMatches)
}
