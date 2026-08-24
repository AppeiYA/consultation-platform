package http

import (
	"github.com/AppeiYA/consultation-platform/internal/shared/logger"
	"github.com/AppeiYA/consultation-platform/internal/shared/response"
	"github.com/AppeiYA/consultation-platform/internal/shared/session"
	"github.com/gofiber/fiber/v2"
	shared_http "github.com/AppeiYA/consultation-platform/internal/shared/adapters/http"
	"go.uber.org/zap"
)

func (h *ConsultantHandler) DeactivateAvailability(c *fiber.Ctx) error {
	claims := c.Locals(session.ContextClaimsKey).(*session.Claims)
	availabilityID := c.Params("availabilityID")
	err := h.ConsultantModule.DeactivateAvailability.Execute(c.Context(), claims.UserID, availabilityID)
	if err != nil {
		logger.Error("Error deactivating availability at handler.DeactivateAvailability", zap.Error(err))
		return response.Error(c, shared_http.StatusFor(err), err.Error(), nil)
	}

	return response.JSON(c, fiber.StatusOK, "Availability deactivated successfully", nil)
}