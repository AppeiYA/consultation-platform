package http

import (
	"github.com/AppeiYA/consultation-platform/internal/consultationcase/adapters/inbound/http/dto"
	"github.com/AppeiYA/consultation-platform/internal/shared/response"
)

var (
	_ = dto.ConsultationCasesDTO{}
	_ = response.Response{}
	_ = response.ErrorResponse{}
)

// GetConsultationCaseByID godoc
// @Summary Get consultation case by ID
// @Description Retrieve details of a consultation case belonging to the authenticated client user.
// @Tags ConsultationCase
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Consultation Case ID"
// @Success 200 {object} response.Response{data=dto.ConsultationCasesDTO} "Case fetched successfully"
// @Failure 400 {object} response.ErrorResponse "Bad request (missing case ID)"
// @Failure 401 {object} response.ErrorResponse "Unauthorized (authentication required or case belongs to another user)"
// @Failure 404 {object} response.ErrorResponse "Case not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /consultation-cases/{id} [get]
func _GetConsultationCaseByID() {}
