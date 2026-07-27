package http

import (
	"github.com/AppeiYA/consultation-platform/internal/identity/adapters/inbound/http/dto"
	shared_http "github.com/AppeiYA/consultation-platform/internal/shared/adapters/http"
	"github.com/AppeiYA/consultation-platform/internal/shared/logger"
	"github.com/AppeiYA/consultation-platform/internal/shared/response"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func (h *IdentityHandler) Register(c *fiber.Ctx) error {
	var req dto.RegisterRequest

	if err := c.BodyParser(&req); err != nil {
		logger.Error("failed to parse request", zap.Error(err))
		return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	params := req.ToRegisterParams()

	res, err := h.identityModule.RegisterUser.Execute(c.Context(), *params)
	if err != nil {
		logger.Error("failed to register user", zap.Error(err))
		return response.Error(c, shared_http.StatusFor(err), err.Error(), nil)
	}

	return response.JSON(c, fiber.StatusCreated, "user registered successfully", dto.NewRegisterUserResponse(res.UserID))
}