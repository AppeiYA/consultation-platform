package mocks

type MockIDGenerator struct {
	GenerateFn func() (string, error)
}

func (m *MockIDGenerator) Generate() (string, error) {
	return m.GenerateFn()
}