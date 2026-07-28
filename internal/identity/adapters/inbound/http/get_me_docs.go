package http

import (
	"github.com/AppeiYA/consultation-platform/internal/identity/adapters/inbound/http/dto"
	"github.com/AppeiYA/consultation-platform/internal/shared/response"
)

var (
	_ = dto.GetMeResponse{}
	_ = response.Response{}
	_ = response.ErrorResponse{}
)

// GetMe godoc
// @Summary Get current user profile
// @Description Retrieve profile details of the authenticated user.
// @Tags Identity
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=dto.GetMeResponse} "User details fetched successfully"
// @Failure 401 {object} response.ErrorResponse "Unauthorized (missing or invalid session token)"
// @Failure 404 {object} response.ErrorResponse "User not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Security SessionAuth
// @Router /identity/me [get]
func _GetMe() {}
