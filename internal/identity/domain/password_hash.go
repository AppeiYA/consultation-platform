package domain

type PasswordHash struct {
	value string
}

func NewPasswordHash(value string) PasswordHash {
	return PasswordHash{
		value: value,
	}
}

func (p PasswordHash) String() string {
	return p.value
}	