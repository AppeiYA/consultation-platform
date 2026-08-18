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

	// Protected
	consultantGroup.Post(
		"/register",
		identityMiddleware.Authenticate,
		consultantHandler.RegisterConsultant,
	)
	consultantGroup.Get(
		"/user",
		identityMiddleware.Authenticate,
		consultantHandler.GetConsultantByUserID,
	)

	// Public
	consultantGroup.Get("/:id", consultantHandler.GetConsultantByID)
	// consultantGroup.Get("/", consultantHandler.ListConsultants)

	consultantGroup.Put(
		"/profile",
		identityMiddleware.Authenticate,
		consultantHandler.UpdateConsultant,
	)
	consultantGroup.Post(
		"/verification",
		identityMiddleware.Authenticate,
		consultantHandler.SubmitVerification,
	)
}