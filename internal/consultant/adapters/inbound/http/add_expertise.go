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

func (h *ConsultantHandler) AddExpertise(c *fiber.Ctx) error {
	claims := c.Locals(session.ContextClaimsKey).(*session.Claims)

	var req dto.AddExpertiseDTO
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid request body", nil)
	}

	res, err := h.ConsultantModule.AddExpertise.Execute(c.Context(), claims.UserID, usecase_dto.AddExpertiseDTO{
		Name: req.Name,
	})
	if err != nil {
		logger.Error("error adding expertise", zap.Error(err), zap.String("user_id", claims.UserID))
		return response.Error(c, shared_http.StatusFor(err), err.Error(), nil)
	}

	return response.JSON(c, fiber.StatusCreated, "Expertise added successfully", res)
}
