package http

import (
	"github.com/AppeiYA/consultation-platform/internal/consultant/adapters/inbound/http/dto"
	"github.com/AppeiYA/consultation-platform/internal/shared/response"
)

var (
	_ = dto.UpdateConsultantModel{}
	_ = response.Response{}
	_ = response.ErrorResponse{}
)

// UpdateConsultant godoc
// @Summary Update consultant profile
// @Description Update the profile details of the authenticated consultant including profession, display name, bio, and years of experience.
// @Tags Consultant
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body dto.UpdateConsultantModel true "Consultant Update Request"
// @Success 200 {object} response.Response "Consultant profile updated successfully"
// @Failure 400 {object} response.ErrorResponse "Bad request (invalid body payload or validation error)"
// @Failure 401 {object} response.ErrorResponse "Unauthorized (authentication required)"
// @Failure 404 {object} response.ErrorResponse "Consultant not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /consultants/profile [put]
func _UpdateConsultant() {}
