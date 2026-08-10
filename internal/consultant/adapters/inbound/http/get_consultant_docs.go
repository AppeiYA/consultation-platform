package http

import (
	"github.com/AppeiYA/consultation-platform/internal/consultant/adapters/inbound/http/dto"
	"github.com/AppeiYA/consultation-platform/internal/shared/response"
)

var (
	_ = dto.PublicConsultantResponseDTO{}
	_ = dto.PrivateConsultantResponseDTO{}
	_ = response.Response{}
	_ = response.ErrorResponse{}
)

// GetConsultantByID godoc
// @Summary Get public consultant profile by ID
// @Description Retrieve public profile details of a consultant by their consultant ID.
// @Tags Consultant
// @Accept json
// @Produce json
// @Param id path string true "Consultant ID"
// @Success 200 {object} response.Response{data=dto.PublicConsultantResponseDTO} "Consultant fetched successfully"
// @Failure 404 {object} response.ErrorResponse "Consultant not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /consultants/{id} [get]
func _GetConsultantByID() {}

// GetConsultantByUserID godoc
// @Summary Get private consultant profile for authenticated user
// @Description Retrieve full consultant profile details of the authenticated user.
// @Tags Consultant
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.Response{data=dto.PrivateConsultantResponseDTO} "Consultant fetched successfully"
// @Failure 401 {object} response.ErrorResponse "Unauthorized (authentication required)"
// @Failure 404 {object} response.ErrorResponse "Consultant not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /consultants/user [get]
func _GetConsultantByUserID() {}
