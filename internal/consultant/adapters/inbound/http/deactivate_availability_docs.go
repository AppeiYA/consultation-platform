package http

import (
	"github.com/AppeiYA/consultation-platform/internal/shared/response"
)

var (
	_ = response.Response{}
	_ = response.ErrorResponse{}
)

// DeactivateAvailability godoc
// @Summary Deactivate consultant availability
// @Description Deactivate an active recurring availability time slot for the authenticated consultant.
// @Tags Consultant
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param availabilityID path string true "Availability ID"
// @Success 200 {object} response.Response "Availability deactivated successfully"
// @Failure 401 {object} response.ErrorResponse "Unauthorized (authentication required)"
// @Failure 404 {object} response.ErrorResponse "Consultant or availability not found"
// @Failure 409 {object} response.ErrorResponse "Conflict (availability is already deactivated)"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /consultants/availability/{availabilityID}/deactivate [patch]
func _DeactivateAvailability() {}
