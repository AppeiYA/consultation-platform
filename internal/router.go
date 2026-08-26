package app

import (
	identity_http "github.com/AppeiYA/consultation-platform/internal/identity/adapters/inbound/http"
	consultant_http "github.com/AppeiYA/consultation-platform/internal/consultant/adapters/inbound/http"
	consultationcase_http "github.com/AppeiYA/consultation-platform/internal/consultationcase/adapters/inbound/http"
	"github.com/AppeiYA/consultation-platform/internal/shared/response"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/contrib/swagger"
)

func SetUpRouter(app *App) {
	app.fiber.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to CONSULTATION API SERVICE")
	})

	swaggerConfig := swagger.Config{
		BasePath: "/api/v1",
		FilePath: "./docs/swagger.json",
		Path:     "swagger",
		Title:    "Consultation Platform API Documentation",
	}

	app.fiber.Use(swagger.New(swaggerConfig))

	v1 := app.fiber.Group("/api/v1")

	v1.Get("/health", func(c *fiber.Ctx) error {
		return response.JSON(c, fiber.StatusOK, "API is Healthy", nil)
	})

	identity_http.SetUpRouter(v1, app.identityHandler, app.identityAuthMiddleware)
	consultant_http.RegisterConsultantRoutes(v1, app.consultantHandler, app.identityAuthMiddleware)
	consultationcase_http.RegisterConsultationCaseRoutes(v1, app.consultationCaseHandler, app.identityAuthMiddleware)

	app.fiber.Use(func(c *fiber.Ctx) error {
		return response.Error(c, fiber.StatusNotFound, "route not found", nil)
	})
}