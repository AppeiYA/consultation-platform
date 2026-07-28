package http

import (
	"github.com/AppeiYA/consultation-platform/internal/identity/adapters/inbound/http/dto"
	"github.com/AppeiYA/consultation-platform/internal/shared/response"
)

var (
	_ = dto.LoginRequest{}
	_ = response.Response{}
	_ = response.ErrorResponse{}
)

// Login godoc
// @Summary User login
// @Description Authenticate a user with email and password and establish a session.
// @Tags Identity
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login Request"
// @Success 200 {object} response.Response "Login successful"
// @Failure 400 {object} response.ErrorResponse "Bad request (invalid body payload)"
// @Failure 401 {object} response.ErrorResponse "Unauthorized (invalid credentials)"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /identity/login [post]
func _Login() {}
