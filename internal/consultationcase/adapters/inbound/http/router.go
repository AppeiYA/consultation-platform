package http

import (
	identity_middleware "github.com/AppeiYA/consultation-platform/internal/identity/adapters/inbound/http/middleware"
	"github.com/gofiber/fiber/v2"
)

func RegisterConsultationCaseRoutes(router fiber.Router, consultationCaseHandler *ConsultationCaseHandler, identityMiddleware *identity_middleware.AuthenticationMiddleware) {
	consultationCaseRouter := router.Group("/consultation-cases")
	consultationCaseRouter.Use(identityMiddleware.Authenticate)

	consultationCaseRouter.Post("/", consultationCaseHandler.CreateConsultationCase)
	consultationCaseRouter.Get("/", consultationCaseHandler.ListConsultationCases)
	consultationCaseRouter.Get("/:id", consultationCaseHandler.GetConsultationCaseByID)
	consultationCaseRouter.Patch("/:id", consultationCaseHandler.UpdateConsultationCase)
	consultationCaseRouter.Put("/:id", consultationCaseHandler.UpdateConsultationCase)
	consultationCaseRouter.Delete("/:id", consultationCaseHandler.DeleteConsultationCase)
}