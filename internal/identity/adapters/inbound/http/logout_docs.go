package http

import (
	"github.com/AppeiYA/consultation-platform/internal/shared/response"
)

var (
	_ = response.Response{}
	_ = response.ErrorResponse{}
)

// Logout godoc
// @Summary Logout user session
// @Description Invalidate the active user session and clear the authentication cookie.
// @Tags Identity
// @Accept json
// @Produce json
// @Success 200 {object} response.Response "Logout successful"
// @Failure 401 {object} response.ErrorResponse "Unauthorized (missing or invalid session cookie)"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Security SessionAuth
// @Router /identity/logout [post]
func _Logout() {}
