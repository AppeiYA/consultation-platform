package http

import (
	"github.com/AppeiYA/consultation-platform/internal/shared/response"
)

var (
	_ = response.Response{}
	_ = response.ErrorResponse{}
)

// RemoveExpertise godoc
// @Summary Remove an expertise from consultant profile
// @Description Delete a specific expertise by its ID from the authenticated consultant profile.
// @Tags Consultant
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param expertiseID path string true "Expertise ID"
// @Success 200 {object} response.Response "Expertise removed successfully"
// @Failure 400 {object} response.ErrorResponse "Bad request (missing expertise ID)"
// @Failure 401 {object} response.ErrorResponse "Unauthorized (authentication required)"
// @Failure 404 {object} response.ErrorResponse "Not found (consultant or expertise not found)"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /consultants/me/expertises/{expertiseID} [delete]
func _RemoveExpertise() {}
