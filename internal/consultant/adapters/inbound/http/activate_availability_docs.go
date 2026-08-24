package http

import (
	"github.com/AppeiYA/consultation-platform/internal/shared/response"
)

var (
	_ = response.Response{}
	_ = response.ErrorResponse{}
)

// ActivateAvailability godoc
// @Summary Activate consultant availability
// @Description Activate a previously deactivated recurring availability time slot for the authenticated consultant.
// @Tags Consultant
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param availabilityID path string true "Availability ID"
// @Success 200 {object} response.Response "Availability activated successfully"
// @Failure 401 {object} response.ErrorResponse "Unauthorized (authentication required)"
// @Failure 404 {object} response.ErrorResponse "Consultant or availability not found"
// @Failure 409 {object} response.ErrorResponse "Conflict (availability is already activated)"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /consultants/availability/{availabilityID}/activate [patch]
func _ActivateAvailability() {}
