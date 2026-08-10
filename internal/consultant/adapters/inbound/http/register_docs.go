package http

import (
	"github.com/AppeiYA/consultation-platform/internal/consultant/adapters/inbound/http/dto"
	"github.com/AppeiYA/consultation-platform/internal/shared/response"
)

var (
	_ = dto.RegisterConsultantDTO{}
	_ = response.Response{}
	_ = response.ErrorResponse{}
)

// RegisterConsultant godoc
// @Summary Register a new consultant profile
// @Description Register a new consultant profile for the authenticated user with profession, display name, bio, and years of experience.
// @Tags Consultant
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body dto.RegisterConsultantDTO true "Consultant Registration Request"
// @Success 201 {object} response.Response "Consultant profile created. Please verify"
// @Failure 400 {object} response.ErrorResponse "Bad request (invalid body payload or validation error)"
// @Failure 401 {object} response.ErrorResponse "Unauthorized (authentication required)"
// @Failure 409 {object} response.ErrorResponse "Conflict (consultant profile already exists for user)"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /consultants/register [post]
func _RegisterConsultant() {}
