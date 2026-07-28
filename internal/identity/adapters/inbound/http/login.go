package http

import (
	"github.com/AppeiYA/consultation-platform/internal/identity/adapters/inbound/http/dto"
	"github.com/AppeiYA/consultation-platform/internal/shared/adapters/http"
	"github.com/AppeiYA/consultation-platform/internal/shared/logger"
	"github.com/AppeiYA/consultation-platform/internal/shared/response"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func (i *IdentityHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Error("failed to parse request", zap.Error(err))
		return response.Error(c, fiber.StatusBadRequest, "failed to parse request", err.Error())
	}

	params := req.ToLoginParams()

	result, err := i.identityModule.LoginUser.Execute(c.Context(), *params)
	if err != nil {
		logger.Error("failed to login", zap.Error(err))
		return response.Error(c, http.StatusFor(err), "failed to login", err.Error())
	}

	i.cookieManager.SetSession(c, result.SessionToken.String())

	return response.JSON(c, fiber.StatusOK, "login success", nil)
}