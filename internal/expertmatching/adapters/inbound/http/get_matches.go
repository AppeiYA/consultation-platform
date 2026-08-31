package http

import (
	"strconv"

	"github.com/AppeiYA/consultation-platform/internal/expertmatching/adapters/inbound/http/dto"
	"github.com/AppeiYA/consultation-platform/internal/expertmatching/ports/inbound"
	shared_http "github.com/AppeiYA/consultation-platform/internal/shared/adapters/http"
	"github.com/AppeiYA/consultation-platform/internal/shared/logger"
	"github.com/AppeiYA/consultation-platform/internal/shared/response"
	"github.com/AppeiYA/consultation-platform/internal/shared/session"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// GetCaseMatches retrieves the latest matching result for a case.
func (h *ExpertMatchingHandler) GetCaseMatches(c *fiber.Ctx) error {
	_ = c.Locals(session.ContextClaimsKey).(*session.Claims)

	caseID := c.Params("id")
	if caseID == "" {
		return response.Error(c, fiber.StatusBadRequest, "case id is required", nil)
	}

	topN := 5
	if topNQuery := c.Query("top_n"); topNQuery != "" {
		if parsed, err := strconv.Atoi(topNQuery); err == nil && parsed > 0 {
			topN = parsed
		}
	}

	result, err := h.expertMatchingModule.GetMatchingResult.Execute(c.Context(), inbound.GetMatchingResultRequest{
		CaseID: caseID,
		TopN:   topN,
	})
	if err != nil {
		logger.Error(
			"error getting matching result at ExpertMatchingHandler.GetCaseMatches",
			zap.Error(err),
			zap.String("case_id", caseID),
		)
		return response.Error(c, shared_http.StatusFor(err), err.Error(), nil)
	}

	resp := dto.FromInboundToMatchingResultResponse(result)
	return response.JSON(c, fiber.StatusOK, "Matching result retrieved successfully", resp)
}
