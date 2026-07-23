package domain

import "time"

type User struct {
	id           string
	firstName    string
	lastName     string
	email         Email
	passwordHash PasswordHash
	role          Role
	isDeleted     bool
	createdAt     time.Time
	updatedAt     time.Time
}

func NewUser(
	id string,
	firstName string,
	lastName string,
	email Email,
	passwordHash PasswordHash,
	role Role,
	now time.Time,
) *User {
	return &User{
		id:           id,
		firstName:    firstName,
		lastName:     lastName,
		email:        email,
		passwordHash: passwordHash,
		role:         role,
		createdAt:    now,
		updatedAt:    now,
	}
}

// Getters

func (u *User) ID() string               { return u.id }
func (u *User) FirstName() string        { return u.firstName }
func (u *User) LastName() string         { return u.lastName }
func (u *User) Email() Email             { return u.email }
func (u *User) PasswordHash() PasswordHash { return u.passwordHash }
func (u *User) Role() Role               { return u.role }
func (u *User) IsDeleted() bool          { return u.isDeleted }
func (u *User) CreatedAt() time.Time     { return u.createdAt }
func (u *User) UpdatedAt() time.Time     { return u.updatedAt }


// setters
func (u *User) ChangeName(firstName, lastName string, now time.Time) {
	u.firstName = firstName
	u.lastName = lastName
	u.updatedAt = now
}

func (u *User) ChangePassword(hash PasswordHash, now time.Time) {
	u.passwordHash = hash
	u.updatedAt = now
}

func (u *User) Delete(now time.Time) {
	u.isDeleted = true
	u.updatedAt = now
}

func (u *User) Restore(now time.Time) {
	u.isDeleted = false
	u.updatedAt = now
}

// func (u *User) PromoteToConsultant(now time.Time) {
// 	u.role = RoleConsultant
// 	u.updatedAt = now
// }

// func (u *User) PromoteToAdmin(now time.Time) {
// 	u.role = RoleAdmin
// 	u.updatedAt = now
// }