package http

import (
	"github.com/AppeiYA/consultation-platform/internal/consultationcase/adapters/inbound/http/dto"
	shared_http "github.com/AppeiYA/consultation-platform/internal/shared/adapters/http"
	"github.com/AppeiYA/consultation-platform/internal/shared/logger"
	"github.com/AppeiYA/consultation-platform/internal/shared/response"
	"github.com/AppeiYA/consultation-platform/internal/shared/session"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func (h *ConsultationCaseHandler) ListConsultationCases(c *fiber.Ctx) error {
	claims := c.Locals(session.ContextClaimsKey).(*session.Claims)
	cases, err := h.consultationCaseModule.ListCase.Execute(c.Context(), claims.UserID)
	if err != nil {
		logger.Error(
			"error listing consultation cases at ConsultationCaseHandler.ListConsultationCases",
			zap.Error(err),
		)
		return response.Error(c, shared_http.StatusFor(err), err.Error(), nil)
	}

	var resp []dto.ConsultationCasesDTO
	for _, consultationCase := range cases {
		resp = append(resp, dto.FromDomainToConsultationCase(consultationCase))
	}
	return response.JSON(c, fiber.StatusOK, "Cases fetched successfully", resp)
}