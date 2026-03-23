package goutils

import "time"

// TimeProvider abstracts access to current time.
type TimeProvider interface {
	Now() time.Time
}

// MockTimeProvider is a mutable TimeProvider used in tests.
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

// SystemTimeProvider returns a TimeProvider backed by time.Now().
func SystemTimeProvider() TimeProvider {
	return &systemTimeProvider{}
}

// NewMockTimeProvider returns a mock provider initialized with current system
// time.
func NewMockTimeProvider() MockTimeProvider {
	return &mockTimeProviderImpl{
		t: time.Now(),
	}
}
