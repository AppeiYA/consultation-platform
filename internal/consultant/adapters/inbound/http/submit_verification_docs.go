package http

import (
	"github.com/AppeiYA/consultation-platform/internal/consultant/adapters/inbound/http/dto"
	"github.com/AppeiYA/consultation-platform/internal/shared/response"
)

var (
	_ = dto.SubmitVerificationResponseDTO{}
	_ = response.Response{}
	_ = response.ErrorResponse{}
)

// SubmitVerification godoc
// @Summary Submit consultant verification inquiry
// @Description Initiate identity and credentials verification inquiry for the authenticated consultant.
// @Tags Consultant
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.Response "Verification submitted successfully"
// @Failure 400 {object} response.ErrorResponse "Bad request (verification already pending/in review/approved, or invalid status)"
// @Failure 401 {object} response.ErrorResponse "Unauthorized (authentication required)"
// @Failure 404 {object} response.ErrorResponse "Consultant not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /consultants/verification [post]
func _SubmitVerification() {}
