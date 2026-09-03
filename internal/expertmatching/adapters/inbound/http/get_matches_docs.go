package http

import (
	"github.com/AppeiYA/consultation-platform/internal/expertmatching/adapters/inbound/http/dto"
	"github.com/AppeiYA/consultation-platform/internal/shared/response"
)

var (
	_ = dto.MatchingResultResponse{}
	_ = response.Response{}
	_ = response.ErrorResponse{}
)

// GetCaseMatches godoc
// @Summary Get matching result for a consultation case
// @Description Retrieves the latest completed matching result and ranked candidate consultants for a consultation case.
// @Tags ExpertMatching
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Consultation Case ID"
// @Param top_n query int false "Top N ranked candidates to return (default: 5)"
// @Success 200 {object} response.Response{data=dto.MatchingResultResponse} "Matching result retrieved successfully"
// @Failure 400 {object} response.ErrorResponse "Bad request (missing case ID)"
// @Failure 401 {object} response.ErrorResponse "Unauthorized (authentication required)"
// @Failure 404 {object} response.ErrorResponse "Not found (matching run not found)"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /expert-matching/cases/{id}/matches [get]
func _GetCaseMatches() {}
