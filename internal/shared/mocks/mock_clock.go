package mocks 

import "time"

type MockClock struct {
	NowFn func() time.Time
}

func (m *MockClock) Now() time.Time {
	return m.NowFn()
}