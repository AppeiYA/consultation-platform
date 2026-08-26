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

func (h *ConsultationCaseHandler) CreateConsultationCase(c *fiber.Ctx) error {
	claims := c.Locals(session.ContextClaimsKey).(*session.Claims)
	var req dto.CreateConsultationCaseDTO
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
	}
	if err := req.Validate(); err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	if err := h.consultationCaseModule.SaveCase.Execute(
		c.Context(), 
		claims.UserID, 
		req.ToUsecaseDTO(),
	); err != nil {
		logger.Error(
			"error creating case at ConsultationCaseHandler.CreateConsultationCase",
			zap.Error(err),
	    )
		return response.Error(c, shared_http.StatusFor(err), err.Error(), nil)
	}

	return response.JSON(c, fiber.StatusCreated, "Case created successfully", nil)
}