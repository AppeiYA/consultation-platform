package http

import (
	"github.com/AppeiYA/consultation-platform/internal/consultant/adapters/inbound/http/dto"
	usecase_dto "github.com/AppeiYA/consultation-platform/internal/consultant/usecase/dto"
	shared_http "github.com/AppeiYA/consultation-platform/internal/shared/adapters/http"
	"github.com/AppeiYA/consultation-platform/internal/shared/logger"
	"github.com/AppeiYA/consultation-platform/internal/shared/response"
	"github.com/AppeiYA/consultation-platform/internal/shared/session"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func (h *ConsultantHandler) ReplaceExpertises(c *fiber.Ctx) error {
	claims := c.Locals(session.ContextClaimsKey).(*session.Claims)

	var req dto.ReplaceExpertisesDTO
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid request body", nil)
	}

	res, err := h.ConsultantModule.ReplaceExpertises.Execute(c.Context(), claims.UserID, usecase_dto.ReplaceExpertisesDTO{
		Expertises: req.Expertises,
	})
	if err != nil {
		logger.Error("error replacing expertises", zap.Error(err), zap.String("user_id", claims.UserID))
		return response.Error(c, shared_http.StatusFor(err), err.Error(), nil)
	}

	return response.JSON(c, fiber.StatusOK, "Expertises updated successfully", res)
}
