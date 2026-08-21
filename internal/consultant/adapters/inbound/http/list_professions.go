package http

import (
	"github.com/AppeiYA/consultation-platform/internal/consultant/adapters/inbound/http/dto"
	shared_http "github.com/AppeiYA/consultation-platform/internal/shared/adapters/http"
	"github.com/AppeiYA/consultation-platform/internal/shared/logger"
	"github.com/AppeiYA/consultation-platform/internal/shared/response"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func (h *ConsultantHandler) ListProfessions(c *fiber.Ctx) error {
	professions, err := h.ConsultantModule.ListProfessions.Execute(c.Context())
	if err != nil {
		logger.Error("Unable to fetch professions", zap.Error(err))
		return response.Error(c, shared_http.StatusFor(err), err.Error(), nil)
	}

	var resp []*dto.ListProfessionsResponse
	for _, profession := range professions {
		resp = append(resp, dto.ProfessionFromDomain(profession))
	}

	return response.JSON(c, fiber.StatusOK, "Professions fetched successfully", resp)
}