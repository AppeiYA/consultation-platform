package http

import (
	"github.com/AppeiYA/consultation-platform/internal/consultant/adapters/inbound/http/dto"
	"github.com/AppeiYA/consultation-platform/internal/shared/response"
)

var (
	_ = dto.ReplaceExpertisesDTO{}
	_ = dto.ExpertiseResponseDTO{}
	_ = response.Response{}
	_ = response.ErrorResponse{}
)

// ReplaceExpertises godoc
// @Summary Replace all expertises of the authenticated consultant
// @Description Replaces the consultant's entire list of expertises in bulk.
// @Tags Consultant
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body dto.ReplaceExpertisesDTO true "Replace Expertises Request"
// @Success 200 {object} response.Response{data=[]dto.ExpertiseResponseDTO} "Expertises updated successfully"
// @Failure 400 {object} response.ErrorResponse "Bad request (invalid payload)"
// @Failure 401 {object} response.ErrorResponse "Unauthorized (authentication required)"
// @Failure 404 {object} response.ErrorResponse "Not found (consultant profile not found for user)"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /consultants/me/expertises [put]
func _ReplaceExpertises() {}
