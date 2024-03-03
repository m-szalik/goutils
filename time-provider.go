package goutils

import "time"

// TimeProvider abstraction for receiving current time
type TimeProvider interface {
	Now() time.Time
}
type MockTimeProvider interface {
	TimeProvider
	Add(delta time.Duration)
}

type mockTimeProviderImpl struct {
	t time.Time
}

func (m *mockTimeProviderImpl) Now() time.Time {
	return m.t
}

func (m *mockTimeProviderImpl) Add(duration time.Duration) {
	m.t = m.t.Add(duration)
}

type systemTimeProvider struct{}

func (s systemTimeProvider) Now() time.Time {
	return time.Now()
}

func SystemTimeProvider() TimeProvider {
	return &systemTimeProvider{}
}

func NewMockTimeProvider() MockTimeProvider {
	return &mockTimeProviderImpl{
		t: time.Now(),
	}
}
