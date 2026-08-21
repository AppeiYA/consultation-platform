package domain

import "time"

const ProfessionIDPrefix = "prof"

type Profession struct {
	id string 
	name string
	createdAt time.Time
}

func NewProfession(id string, name string, createdAt time.Time) Profession {
	return Profession{
		id: id,
		name: name,
		createdAt: createdAt,
	}
}

func ReconstituteProfession(id string, name string, createdAt time.Time) Profession {
	return Profession{
		id: id,
		name: name,
		createdAt: createdAt,
	}
}

func (p Profession) ID() string {
	return p.id
}

func (p Profession) Name() string {
	return p.name
}

func (p Profession) CreatedAt() time.Time {
	return p.createdAt
}

// type Profession string

// const (
// 	SoftwareEngineer Profession = "SOFTWARE_ENGINEER"
// 	Lawyer           Profession = "LAWYER"
// 	Doctor           Profession = "DOCTOR"
// 	Accountant       Profession = "ACCOUNTANT"
// 	Therapist Profession = "THERAPIST"
// 	Clergy Profession = "CLERGY"
// )

// func NewProfession(profession string) (Profession, error) {
// 	switch Profession(profession) {
// 	case SoftwareEngineer, Lawyer, Doctor, Accountant, Therapist, Clergy:
// 		return Profession(profession), nil
// 	default:
// 		return Profession(""), ErrInvalidProfession
// 	}
// }