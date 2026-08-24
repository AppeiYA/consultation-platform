package dto

import "github.com/AppeiYA/consultation-platform/internal/consultant/domain"

type GetAvailabilityResponse struct {
	AvailabilityID string `json:"availability_id"`
	IsActive       bool   `json:"is_active"`
	DayOfWeek string `json:"day_of_week"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

func AvailabilityFromDomain(availability *domain.ConsultantAvailability) *GetAvailabilityResponse {
	return &GetAvailabilityResponse{
		AvailabilityID: availability.ID(),
		IsActive:       availability.IsActive(),
		DayOfWeek: availability.DayOfWeek().String(),
		StartTime: availability.StartTime().String(),
		EndTime:   availability.EndTime().String(),
	}
}