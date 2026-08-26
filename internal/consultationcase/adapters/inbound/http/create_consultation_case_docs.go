package http

import (
	"github.com/AppeiYA/consultation-platform/internal/consultationcase/adapters/inbound/http/dto"
	"github.com/AppeiYA/consultation-platform/internal/shared/response"
)

var (
	_ = dto.CreateConsultationCaseDTO{}
	_ = response.Response{}
	_ = response.ErrorResponse{}
)

// CreateConsultationCase godoc
// @Summary Create a consultation case
// @Description Create a new consultation case for the authenticated client user.
// @Tags ConsultationCase
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body dto.CreateConsultationCaseDTO true "Create Consultation Case Request"
// @Success 201 {object} response.Response "Case created successfully"
// @Failure 400 {object} response.ErrorResponse "Bad request (invalid payload or domain validation failure)"
// @Failure 401 {object} response.ErrorResponse "Unauthorized (authentication required)"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /consultation-cases [post]
func _CreateConsultationCase() {}
