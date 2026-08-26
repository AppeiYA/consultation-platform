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

func (h *ConsultationCaseHandler) UpdateConsultationCase(c *fiber.Ctx) error {
	claims := c.Locals(session.ContextClaimsKey).(*session.Claims)
	caseID := c.Params("id")
	if caseID == "" {
		return response.Error(c, fiber.StatusBadRequest, "case id is required", nil)
	}

	var req dto.UpdateConsultationCaseDTO
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid request body", nil)
	}

	if err := req.Validate(); err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	if err := h.consultationCaseModule.UpdateCase.Execute(
		c.Context(),
		claims.UserID,
		caseID,
		req.ToUsecaseDTO(),
	); err != nil {
		logger.Error(
			"error updating consultation case at ConsultationCaseHandler.UpdateConsultationCase",
			zap.Error(err),
		)
		return response.Error(c, shared_http.StatusFor(err), err.Error(), nil)
	}

	return response.JSON(c, fiber.StatusOK, "Case updated successfully", nil)
}
