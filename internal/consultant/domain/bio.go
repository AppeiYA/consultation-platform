package domain

type Bio struct {
	value string
}

func NewBio(bio string) (Bio, error) {
	if len(bio) > 1000 {
		return Bio{}, ErrBioOverflow
	}

	if len(bio) == 0 {
		return Bio{}, ErrEmptyBio
	}

	return Bio{value: bio}, nil
}