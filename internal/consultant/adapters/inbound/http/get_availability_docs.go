package http

import (
	"github.com/AppeiYA/consultation-platform/internal/consultant/adapters/inbound/http/dto"
	"github.com/AppeiYA/consultation-platform/internal/shared/response"
)

var (
	_ = dto.GetAvailabilityResponse{}
	_ = response.Response{}
	_ = response.ErrorResponse{}
)

// GetAvailability godoc
// @Summary Get consultant availability
// @Description Retrieve weekly recurring availabilities for a specific consultant by their consultant ID.
// @Tags Consultant
// @Accept json
// @Produce json
// @Param consultantID path string true "Consultant ID"
// @Success 200 {object} response.Response{data=[]dto.GetAvailabilityResponse} "Availabilities fetched successfully"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /consultants/{consultantID}/availability [get]
func _GetAvailability() {}
