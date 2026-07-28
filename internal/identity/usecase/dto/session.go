package dto

type SessionClaims struct {
	UserID string
	Email  string
	Role   string
}

type ValidateSessionRequest struct {
	SessionToken string
}

type ValidateSessionResponse struct {
	SessionClaims SessionClaims
}