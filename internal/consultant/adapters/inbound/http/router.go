package http

import (
	identity_middleware "github.com/AppeiYA/consultation-platform/internal/identity/adapters/inbound/http/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterConsultantRoutes(
	router fiber.Router,
	consultantHandler *ConsultantHandler,
	identityMiddleware *identity_middleware.AuthenticationMiddleware,
) {
	consultantGroup := router.Group("/consultants")

	// Public
	// consultantGroup.Get("/:id", consultantHandler.GetConsultant)
	// consultantGroup.Get("/", consultantHandler.ListConsultants)

	// Protected
	consultantGroup.Post(
		"/register",
		identityMiddleware.Authenticate,
		consultantHandler.RegisterConsultant,
	)

	// consultantGroup.Put(
	// 	"/profile",
	// 	identityMiddleware.Authenticate,
	// 	consultantHandler.UpdateConsultant,
	// )
}