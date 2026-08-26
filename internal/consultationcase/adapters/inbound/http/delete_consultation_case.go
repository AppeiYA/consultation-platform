package http

import (
	shared_http "github.com/AppeiYA/consultation-platform/internal/shared/adapters/http"
	"github.com/AppeiYA/consultation-platform/internal/shared/logger"
	"github.com/AppeiYA/consultation-platform/internal/shared/response"
	"github.com/AppeiYA/consultation-platform/internal/shared/session"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func (h *ConsultationCaseHandler) DeleteConsultationCase(c *fiber.Ctx) error {
	claims := c.Locals(session.ContextClaimsKey).(*session.Claims)
	caseID := c.Params("id")
	if caseID == "" {
		return response.Error(c, fiber.StatusBadRequest, "case id is required", nil)
	}

	if err := h.consultationCaseModule.DeleteCase.Execute(
		c.Context(),
		claims.UserID,
		caseID,
	); err != nil {
		logger.Error(
			"error deleting consultation case at ConsultationCaseHandler.DeleteConsultationCase",
			zap.Error(err),
		)
		return response.Error(c, shared_http.StatusFor(err), err.Error(), nil)
	}

	return response.JSON(c, fiber.StatusOK, "Case deleted successfully", nil)
}
