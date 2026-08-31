package http

import (
	shared_http "github.com/AppeiYA/consultation-platform/internal/shared/adapters/http"
	"github.com/AppeiYA/consultation-platform/internal/shared/logger"
	"github.com/AppeiYA/consultation-platform/internal/shared/response"
	"github.com/AppeiYA/consultation-platform/internal/shared/session"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func (h *ConsultantHandler) RemoveExpertise(c *fiber.Ctx) error {
	claims := c.Locals(session.ContextClaimsKey).(*session.Claims)
	expertiseID := c.Params("expertiseID")
	if expertiseID == "" {
		return response.Error(c, fiber.StatusBadRequest, "expertise id is required", nil)
	}

	if err := h.ConsultantModule.RemoveExpertise.Execute(c.Context(), claims.UserID, expertiseID); err != nil {
		logger.Error("error removing expertise", zap.Error(err), zap.String("expertise_id", expertiseID))
		return response.Error(c, shared_http.StatusFor(err), err.Error(), nil)
	}

	return response.JSON(c, fiber.StatusOK, "Expertise removed successfully", nil)
}
