package dto

type GetCurrentUserRequest struct {
	UserID string
}

type GetCurrentUserResponse struct {
	ID        string
	FirstName string
	LastName  string
	Email     string
	Role      string
}