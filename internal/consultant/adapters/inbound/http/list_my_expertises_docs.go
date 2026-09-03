package http

import (
	"github.com/AppeiYA/consultation-platform/internal/consultant/adapters/inbound/http/dto"
	"github.com/AppeiYA/consultation-platform/internal/shared/response"
)

var (
	_ = dto.ExpertiseResponseDTO{}
	_ = response.Response{}
	_ = response.ErrorResponse{}
)

// ListMyExpertises godoc
// @Summary List expertises of the authenticated consultant
// @Description Retrieve all expertise skills and specialties belonging to the authenticated consultant.
// @Tags Consultant
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.Response{data=[]dto.ExpertiseResponseDTO} "Expertises retrieved successfully"
// @Failure 401 {object} response.ErrorResponse "Unauthorized (authentication required)"
// @Failure 404 {object} response.ErrorResponse "Not found (consultant profile not found for user)"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /consultants/me/expertises [get]
func _ListMyExpertises() {}
