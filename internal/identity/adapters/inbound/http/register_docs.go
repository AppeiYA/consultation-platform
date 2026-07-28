package http

import (
	"github.com/AppeiYA/consultation-platform/internal/identity/adapters/inbound/http/dto"
	"github.com/AppeiYA/consultation-platform/internal/shared/response"
)

var (
	_ = dto.RegisterRequest{}
	_ = dto.RegisterUserResponse{}
	_ = response.Response{}
	_ = response.ErrorResponse{}
)

// Register godoc
// @Summary Register a new user
// @Description Register a new user account with first name, last name, email, password, and role.
// @Tags Identity
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "User Registration Request"
// @Success 201 {object} response.Response{data=dto.RegisterUserResponse} "User registered successfully"
// @Failure 400 {object} response.ErrorResponse "Bad request (invalid body payload or validation error)"
// @Failure 409 {object} response.ErrorResponse "Conflict (user with email already exists)"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /identity/register [post]
func _Register() {}