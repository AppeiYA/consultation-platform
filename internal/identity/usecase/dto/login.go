package dto

import "github.com/AppeiYA/consultation-platform/internal/identity/domain"

type LoginRequest struct {
	Email    string
	Password string
}

type LoginResponse struct {
	SessionToken domain.SessionToken
}