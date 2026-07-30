package http

import (
	usecase_dto "github.com/AppeiYA/consultation-platform/internal/identity/usecase/dto"
	shared_http "github.com/AppeiYA/consultation-platform/internal/shared/adapters/http"
	"github.com/AppeiYA/consultation-platform/internal/shared/response"
	"github.com/gofiber/fiber/v2"
)

func (h *IdentityHandler) Logout(c *fiber.Ctx) error {
	sessionToken, exists := h.cookieManager.GetSession(c)
	if !exists {
		return fiber.ErrUnauthorized
	}

	_, err := h.identityModule.LogoutUser.Execute(c.Context(), usecase_dto.LogoutRequest{
		SessionToken: sessionToken,
	})
	if err != nil {
		return response.Error(c, shared_http.StatusFor(err), "failed to logout", err.Error())
	}

	h.cookieManager.DeleteSession(c)

	return response.JSON(c, fiber.StatusOK, "Logout Successful", nil)
}