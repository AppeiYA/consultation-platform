package domain

import "time"

type UserRegistered struct {
	UserID string
	At     time.Time
}

type UserLoggedIn struct {
	UserID string
	At     time.Time
}

type UserLoggedOut struct {
	UserID string
	At     time.Time
}