package session

const ContextClaimsKey = "session.claims"

type Claims struct {
	UserID string
	Email  string
	Role   string
}