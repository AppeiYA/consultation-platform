package http

import (
	"github.com/AppeiYA/consultation-platform/internal/consultant/adapters/inbound/http/dto"
	"github.com/AppeiYA/consultation-platform/internal/shared/response"
)

var (
	_ = dto.CreateAvailabilityRequest{}
	_ = response.Response{}
	_ = response.ErrorResponse{}
)

// CreateAvailability godoc
// @Summary Create consultant availability
// @Description Create a recurring weekly availability time slot for the authenticated consultant.
// @Tags Consultant
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body dto.CreateAvailabilityRequest true "Create Availability Request"
// @Success 201 {object} response.Response "Availability created successfully"
// @Failure 400 {object} response.ErrorResponse "Bad request (invalid body payload, invalid time format, or invalid time range)"
// @Failure 401 {object} response.ErrorResponse "Unauthorized (authentication required)"
// @Failure 404 {object} response.ErrorResponse "Consultant not found"
// @Failure 409 {object} response.ErrorResponse "Conflict (availability overlaps with existing availability)"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /consultants/availability [post]
func _CreateAvailability() {}
