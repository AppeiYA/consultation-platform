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

func (h *ConsultantHandler) SubmitVerification(c *fiber.Ctx) error {
	claims := c.Locals(session.ContextClaimsKey).(*session.Claims)

	res, err := h.ConsultantModule.SubmitVerification.Execute(c.Context(), claims.UserID)
	if err != nil {
		logger.Error("Error submitting verification at handler.SubmitVerification", zap.Error(err))
		return response.Error(c, shared_http.StatusFor(err), err.Error(), nil)
	}

	respDTO := dto.SubmitVerificationResponseDTO{
		VerificationID:    res.VerificationID,
		ProviderReference: res.ProviderReference,
		VerificationURL:   res.VerificationURL,
		Status:            res.Status,
	}

	return response.JSON(c, fiber.StatusOK, "Verification submitted successfully", respDTO)
}
