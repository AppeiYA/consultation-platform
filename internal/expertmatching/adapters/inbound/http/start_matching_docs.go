package http

import (
	"github.com/AppeiYA/consultation-platform/internal/expertmatching/adapters/inbound/http/dto"
	"github.com/AppeiYA/consultation-platform/internal/shared/response"
)

var (
	_ = dto.StartMatchingResponse{}
	_ = response.Response{}
	_ = response.ErrorResponse{}
)

// StartMatching godoc
// @Summary Start matching run for a consultation case
// @Description Initiates an asynchronous expert matching pipeline for the specified consultation case.
// @Tags ExpertMatching
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Consultation Case ID"
// @Success 202 {object} response.Response{data=dto.StartMatchingResponse} "Matching run initiated"
// @Failure 400 {object} response.ErrorResponse "Bad request (missing or invalid case ID)"
// @Failure 401 {object} response.ErrorResponse "Unauthorized (authentication required)"
// @Failure 404 {object} response.ErrorResponse "Not found (consultation case not found)"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /expert-matching/cases/{id}/match [post]
func _StartMatching() {}
