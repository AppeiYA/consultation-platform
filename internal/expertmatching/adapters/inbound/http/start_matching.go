package http

import (
	"github.com/AppeiYA/consultation-platform/internal/expertmatching/adapters/inbound/http/dto"
	shared_http "github.com/AppeiYA/consultation-platform/internal/shared/adapters/http"
	"github.com/AppeiYA/consultation-platform/internal/shared/logger"
	"github.com/AppeiYA/consultation-platform/internal/shared/response"
	"github.com/AppeiYA/consultation-platform/internal/shared/session"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// StartMatching triggers an asynchronous matching run for the given consultation case.
func (h *ExpertMatchingHandler) StartMatching(c *fiber.Ctx) error {
	_ = c.Locals(session.ContextClaimsKey).(*session.Claims)

	caseID := c.Params("id")
	if caseID == "" {
		return response.Error(c, fiber.StatusBadRequest, "case id is required", nil)
	}

	run, err := h.expertMatchingModule.StartMatching.Execute(c.Context(), caseID)
	if err != nil {
		logger.Error(
			"error starting matching run at ExpertMatchingHandler.StartMatching",
			zap.Error(err),
			zap.String("case_id", caseID),
		)
		return response.Error(c, shared_http.StatusFor(err), err.Error(), nil)
	}

	resp := dto.FromDomainToStartMatchingResponse(run)
	return response.JSON(c, fiber.StatusAccepted, "Matching run initiated", resp)
}
