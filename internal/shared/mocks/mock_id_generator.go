package mocks

type MockIDGenerator struct {
	GenerateFn func(prefix string) (string, error)
}

func (m *MockIDGenerator) Generate(prefix string) (string, error) {
	return m.GenerateFn(prefix)
}