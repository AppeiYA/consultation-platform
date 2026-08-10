package http

import (
	"github.com/AppeiYA/consultation-platform/internal/consultant/adapters/inbound/http/dto"
	"github.com/AppeiYA/consultation-platform/internal/shared/logger"
	shared_http "github.com/AppeiYA/consultation-platform/internal/shared/adapters/http"
	"github.com/AppeiYA/consultation-platform/internal/shared/response"
	"github.com/AppeiYA/consultation-platform/internal/shared/session"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func(h *ConsultantHandler) UpdateConsultant(c *fiber.Ctx) error {
	claims := c.Locals(session.ContextClaimsKey).(*session.Claims)
	var req dto.UpdateConsultantModel
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	if err := h.ConsultantModule.UpdateConsultant.Execute(c.Context(), claims.UserID, req.ToUsecaseDTO()); err != nil {
		logger.Error("Error updating consultant at handler.UpdateConsultant", zap.Error(err))
		return response.Error(c, shared_http.StatusFor(err), err.Error(), nil)
	}

	return response.JSON(c, fiber.StatusOK, "Consultant profile updated successfully", nil)
}