package http

import (
	"github.com/AppeiYA/consultation-platform/internal/consultant/adapters/inbound/http/dto"
	shared_http "github.com/AppeiYA/consultation-platform/internal/shared/adapters/http"
	"github.com/AppeiYA/consultation-platform/internal/shared/logger"
	"github.com/AppeiYA/consultation-platform/internal/shared/response"
	"github.com/AppeiYA/consultation-platform/internal/shared/session"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func (r *ConsultantHandler) UpdateAvailability(c *fiber.Ctx) error {
	claims := c.Locals(session.ContextClaimsKey).(*session.Claims)
	var req dto.UpdateAvailabilityRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	if err := req.Validate(); err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	resp, err := r.ConsultantModule.UpdateAvailability.Execute(c.Context(), claims.UserID, req.ToUsecaseDTO())
	if err != nil {
		logger.Error("unable to update availability at ConsultantHandler.UpdateAvailability", zap.Error(err))
		return response.Error(c, shared_http.StatusFor(err), err.Error(), nil)
	}

	return response.JSON(c, fiber.StatusOK, "Availability updated successfully", dto.AvailabilityFromDomain(resp))
}