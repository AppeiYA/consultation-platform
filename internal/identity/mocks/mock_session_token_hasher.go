package mocks 

type MockSessionTokenHasher struct {
	HashFn func(token string) (string, error)
}

func (m *MockSessionTokenHasher) Hash(token string) (string, error) {
	if m.HashFn != nil {
		return m.HashFn(token)
	}
	return "", nil
}