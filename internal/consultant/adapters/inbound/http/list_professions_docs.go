package http

import (
	"github.com/AppeiYA/consultation-platform/internal/consultant/adapters/inbound/http/dto"
	"github.com/AppeiYA/consultation-platform/internal/shared/response"
)

var (
	_ = dto.ListProfessionsResponse{}
	_ = response.Response{}
	_ = response.ErrorResponse{}
)

// ListProfessions godoc
// @Summary List all professions
// @Description Retrieve a list of all active professions available in the system.
// @Tags Consultant
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]dto.ListProfessionsResponse} "Professions fetched successfully"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /consultants/professions [get]
func _ListProfessions() {}
