package domain

type Profession string

const (
	SoftwareEngineer Profession = "SOFTWARE_ENGINEER"
	Lawyer           Profession = "LAWYER"
	Doctor           Profession = "DOCTOR"
	Accountant       Profession = "ACCOUNTANT"
	Therapist Profession = "THERAPIST"
	Clergy Profession = "CLERGY"
)

func NewProfession(profession string) (Profession, error) {
	switch Profession(profession) {
	case SoftwareEngineer, Lawyer, Doctor, Accountant, Therapist, Clergy:
		return Profession(profession), nil
	default:
		return Profession(""), ErrInvalidProfession
	}
}