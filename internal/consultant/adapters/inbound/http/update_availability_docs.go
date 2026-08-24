package http

import (
	"github.com/AppeiYA/consultation-platform/internal/consultant/adapters/inbound/http/dto"
	"github.com/AppeiYA/consultation-platform/internal/shared/response"
)

var (
	_ = dto.UpdateAvailabilityRequest{}
	_ = dto.GetAvailabilityResponse{}
	_ = response.Response{}
	_ = response.ErrorResponse{}
)

// UpdateAvailability godoc
// @Summary Update consultant availability
// @Description Update an existing weekly recurring availability time slot for the authenticated consultant.
// @Tags Consultant
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body dto.UpdateAvailabilityRequest true "Update Availability Request"
// @Success 200 {object} response.Response{data=dto.GetAvailabilityResponse} "Availability updated successfully"
// @Failure 400 {object} response.ErrorResponse "Bad request (invalid body payload, invalid time format, or invalid time range)"
// @Failure 401 {object} response.ErrorResponse "Unauthorized (authentication required)"
// @Failure 404 {object} response.ErrorResponse "Consultant or availability not found"
// @Failure 409 {object} response.ErrorResponse "Conflict (availability overlaps with existing availability)"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /consultants/availability [put]
func _UpdateAvailability() {}
