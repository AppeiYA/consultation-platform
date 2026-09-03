package http

import (
	"github.com/AppeiYA/consultation-platform/internal/consultant/adapters/inbound/http/dto"
	"github.com/AppeiYA/consultation-platform/internal/shared/response"
)

var (
	_ = dto.AddExpertiseDTO{}
	_ = dto.ExpertiseResponseDTO{}
	_ = response.Response{}
	_ = response.ErrorResponse{}
)

// AddExpertise godoc
// @Summary Add a new expertise to consultant profile
// @Description Add a new expertise skill or specialty to the authenticated consultant profile.
// @Tags Consultant
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body dto.AddExpertiseDTO true "Add Expertise Request"
// @Success 201 {object} response.Response{data=dto.ExpertiseResponseDTO} "Expertise added successfully"
// @Failure 400 {object} response.ErrorResponse "Bad request (invalid body payload or empty name)"
// @Failure 401 {object} response.ErrorResponse "Unauthorized (authentication required)"
// @Failure 404 {object} response.ErrorResponse "Not found (consultant profile not found for user)"
// @Failure 409 {object} response.ErrorResponse "Conflict (expertise already exists for consultant)"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /consultants/me/expertises [post]
func _AddExpertise() {}
