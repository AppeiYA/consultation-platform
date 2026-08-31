package http

import (
	shared_http "github.com/AppeiYA/consultation-platform/internal/shared/adapters/http"
	"github.com/AppeiYA/consultation-platform/internal/shared/logger"
	"github.com/AppeiYA/consultation-platform/internal/shared/response"
	"github.com/AppeiYA/consultation-platform/internal/shared/session"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func (h *ConsultantHandler) ListMyExpertises(c *fiber.Ctx) error {
	claims := c.Locals(session.ContextClaimsKey).(*session.Claims)

	res, err := h.ConsultantModule.ListMyExpertises.Execute(c.Context(), claims.UserID)
	if err != nil {
		logger.Error("error listing expertises", zap.Error(err), zap.String("user_id", claims.UserID))
		return response.Error(c, shared_http.StatusFor(err), err.Error(), nil)
	}

	return response.JSON(c, fiber.StatusOK, "Expertises retrieved successfully", res)
}
