package domain

type SessionClaims struct {
	userID string
	email  string
	role   string
}

func (s *SessionClaims) UserID() string {
	return s.userID
}

func (s *SessionClaims) Email() string {
	return s.email
}

func (s *SessionClaims) Role() string {
	return s.role
}

func NewSessionClaimsFromParams(userID, email, role string) *SessionClaims {
	return &SessionClaims{
		userID: userID,
		email:  email,
		role:   role,
	}
}

func NewSessionClaims(session *Session) *SessionClaims {
	return &SessionClaims{
		userID: session.userID,
		email:  session.email,
		role:   session.role,
	}
}
