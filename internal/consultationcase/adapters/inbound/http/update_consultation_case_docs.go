package http

import (
	"github.com/AppeiYA/consultation-platform/internal/consultationcase/adapters/inbound/http/dto"
	"github.com/AppeiYA/consultation-platform/internal/shared/response"
)

var (
	_ = dto.UpdateConsultationCaseDTO{}
	_ = response.Response{}
	_ = response.ErrorResponse{}
)

// UpdateConsultationCase godoc
// @Summary Update consultation case
// @Description Partially update fields of an existing consultation case belonging to the authenticated client.
// @Tags ConsultationCase
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Consultation Case ID"
// @Param request body dto.UpdateConsultationCaseDTO true "Update Consultation Case Request"
// @Success 200 {object} response.Response "Case updated successfully"
// @Failure 400 {object} response.ErrorResponse "Bad request (invalid body or validation failure)"
// @Failure 401 {object} response.ErrorResponse "Unauthorized (authentication required or case belongs to another user)"
// @Failure 404 {object} response.ErrorResponse "Case not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /consultation-cases/{id} [patch]
func _UpdateConsultationCase() {}
