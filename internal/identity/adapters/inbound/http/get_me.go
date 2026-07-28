package http

import (
	"github.com/AppeiYA/consultation-platform/internal/identity/adapters/inbound/http/dto"
	"github.com/AppeiYA/consultation-platform/internal/identity/adapters/inbound/http/middleware"
	usecase_dto "github.com/AppeiYA/consultation-platform/internal/identity/usecase/dto"
	"github.com/AppeiYA/consultation-platform/internal/shared/adapters/http"
	"github.com/AppeiYA/consultation-platform/internal/shared/logger"
	"github.com/AppeiYA/consultation-platform/internal/shared/response"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func (i *IdentityHandler) GetMe(c *fiber.Ctx) error {
	claims := c.Locals(middleware.ContextClaimKey).(*usecase_dto.SessionClaims)

	resp, err := i.identityModule.GetCurrentUser.Execute(c.UserContext(), usecase_dto.GetCurrentUserRequest{
		UserID: claims.UserID,
	})
	if err != nil {
		logger.Error("error getting user: ", zap.Error(err))
		return response.Error(c, http.StatusFor(err), err.Error(), nil)
	}

	result := dto.NewGetMeResponse(&resp)

	return response.JSON(c, fiber.StatusOK, "user fetched successfully", result)
}