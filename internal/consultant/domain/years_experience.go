package domain

type YearsExperience struct {
	value int
}

func NewYearsExperience(years int) (YearsExperience, error) {
	if years <= 0 || years >= 70 {
		return YearsExperience{}, ErrInvalidYearsExperience
	}

	return YearsExperience{value: years}, nil
}

func (y YearsExperience) Int() int {
	return y.value
}