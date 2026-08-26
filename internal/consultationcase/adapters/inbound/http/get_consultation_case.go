package http

import (
	"github.com/AppeiYA/consultation-platform/internal/consultationcase/adapters/inbound/http/dto"
	shared_http "github.com/AppeiYA/consultation-platform/internal/shared/adapters/http"
	"github.com/AppeiYA/consultation-platform/internal/shared/logger"
	"github.com/AppeiYA/consultation-platform/internal/shared/response"
	"github.com/AppeiYA/consultation-platform/internal/shared/session"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func (h *ConsultationCaseHandler) GetConsultationCaseByID(c *fiber.Ctx) error {
	claims := c.Locals(session.ContextClaimsKey).(*session.Claims)
	caseID := c.Params("id")
	if caseID == "" {
		return response.Error(c, fiber.StatusBadRequest, "case id is required", nil)
	}

	consultationCase, err := h.consultationCaseModule.GetCase.Execute(c.Context(), claims.UserID, caseID)
	if err != nil {
		logger.Error(
			"error getting consultation case at ConsultationCaseHandler.GetConsultationCaseByID",
			zap.Error(err),
		)
		return response.Error(c, shared_http.StatusFor(err), err.Error(), nil)
	}

	resp := dto.FromDomainToConsultationCase(consultationCase)
	return response.JSON(c, fiber.StatusOK, "Case fetched successfully", resp)
}
