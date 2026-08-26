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

// ListConsultationCases godoc
// @Summary List client's consultation cases
// @Description Retrieve all consultation cases created by the authenticated client user.
// @Tags ConsultationCase
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.Response{data=[]dto.ConsultationCasesDTO} "Cases fetched successfully"
// @Failure 401 {object} response.ErrorResponse "Unauthorized (authentication required)"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /consultation-cases [get]
func _ListConsultationCases() {}
