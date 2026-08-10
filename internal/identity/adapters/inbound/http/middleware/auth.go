package middleware

import (
	"github.com/AppeiYA/consultation-platform/internal/identity"
	"github.com/AppeiYA/consultation-platform/internal/identity/usecase/dto"
	shared_http "github.com/AppeiYA/consultation-platform/internal/shared/adapters/http"
	"github.com/AppeiYA/consultation-platform/internal/shared/logger"
	"github.com/AppeiYA/consultation-platform/internal/shared/response"
	"github.com/AppeiYA/consultation-platform/internal/shared/session"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type AuthenticationMiddleware struct {
	identityModule *identity.Module
	cookies shared_http.CookieManagerInt
}

func NewAuthenticationMiddleware(
	identityModule *identity.Module,
	cookies shared_http.CookieManagerInt,
) *AuthenticationMiddleware {
	return &AuthenticationMiddleware{
		identityModule: identityModule,
		cookies: cookies,
	}
}

func (m *AuthenticationMiddleware) Authenticate(c *fiber.Ctx) error {
	token, ok := m.cookies.GetSession(c)
	if !ok {
		logger.Error("Unauthorized user", zap.String("error", "unauthorized user"))
		return response.Error(c, fiber.StatusUnauthorized, "User is unauthorized", nil)
	}

	params := dto.ValidateSessionRequest{
		SessionToken: token,
	}

	claims, err := m.identityModule.ValidateSession.Execute(c.Context(), params)
	if err != nil {
		logger.Error("error validating session: ", zap.Error(err))
		return response.Error(c, shared_http.StatusFor(err), "User is unauthorized", nil)
	}

	c.Locals(session.ContextClaimsKey, &session.Claims{
		UserID: claims.SessionClaims.UserID,
		Email:  claims.SessionClaims.Email,
		Role:   claims.SessionClaims.Role,
	})

	return c.Next()	
}