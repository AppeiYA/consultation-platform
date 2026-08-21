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

func (h *ConsultantHandler) GetConsultantByID(c *fiber.Ctx) error {
	consultantID := c.Params("id")
	var getConsultantResp dto.PublicConsultantResponseDTO

	consultant, err := h.ConsultantModule.GetConsultant.ByID(c.Context(), consultantID)
	if err != nil {
		logger.Error("Error getting consultant by ID at handler.GetConsultantByID", zap.Error(err))
		return response.Error(c, shared_http.StatusFor(err), err.Error(), nil)
	}

	getConsultantResp.FromUsecaseDTO(consultant)

	return response.JSON(c, fiber.StatusOK, "Consultant fetched successfully", getConsultantResp)
}

func (h *ConsultantHandler) GetConsultantByUserID(c *fiber.Ctx) error {
	claims := c.Locals(session.ContextClaimsKey).(*session.Claims)
	var getConsultantResp dto.PrivateConsultantResponseDTO

	consultant, err := h.ConsultantModule.GetConsultant.ByUserID(c.Context(), claims.UserID)
	if err != nil {
		logger.Error("Error getting consultant by User ID at handler.GetConsultantByUserID", zap.Error(err))
		return response.Error(c, shared_http.StatusFor(err), err.Error(), nil)
	}

	getConsultantResp.FromUsecaseDTO(consultant)

	return response.JSON(c, fiber.StatusOK, "Consultant fetched successfully", getConsultantResp)
}