package http

import (
	"github.com/AppeiYA/consultation-platform/internal/consultant/adapters/inbound/http/dto"
	shared_http "github.com/AppeiYA/consultation-platform/internal/shared/adapters/http"
	"github.com/AppeiYA/consultation-platform/internal/shared/logger"
	"github.com/AppeiYA/consultation-platform/internal/shared/response"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func (h *ConsultantHandler) GetAvailability(c *fiber.Ctx) error {
	consultantID := c.Params("consultantID")
	availabilities, err := h.ConsultantModule.GetAvailability.Execute(c.Context(), consultantID)
	if err != nil {
		logger.Error("Error getting availability at handler.GetAvailability", zap.Error(err))
		return response.Error(c, shared_http.StatusFor(err), err.Error(), nil)
	}

	resp := make([]*dto.GetAvailabilityResponse, 0, len(availabilities))
	for _, availability := range availabilities {
		resp = append(resp, dto.AvailabilityFromDomain(availability))
	}

	return response.JSON(c, fiber.StatusOK, "Availabilities fetched successfully", resp)
}