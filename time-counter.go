package goutils

import (
	"sync"
	"time"
)

// TimeCounter measures cumulative elapsed time between Start and Stop calls.
type TimeCounter struct {
	timeProvider TimeProvider
	sessionStart *time.Time
	continued    time.Duration
	lock         sync.Mutex
}

// Reset returns the current elapsed duration and resets it to zero.
// If the counter is running, it continues running from the reset point.
func (d *TimeCounter) Reset() time.Duration {
	d.lock.Lock()
	defer d.lock.Unlock()
	now := d.timeProvider.Now()
	if d.sessionStart != nil {
		d.continued += now.Sub(*d.sessionStart)
		d.sessionStart = &now
	}
	dur := d.continued
	d.continued = time.Duration(0)
	return dur
}

// Value returns current elapsed duration including the active session.
func (d *TimeCounter) Value() time.Duration {
	d.lock.Lock()
	defer d.lock.Unlock()
	dur := d.continued
	if d.sessionStart != nil {
		additional := d.timeProvider.Now().Sub(*d.sessionStart)
		dur += additional
	}
	return dur
}

// Start starts measuring time when it is not already running.
func (d *TimeCounter) Start() {
	d.lock.Lock()
	defer d.lock.Unlock()
	if d.sessionStart == nil {
		now := d.timeProvider.Now()
		d.sessionStart = &now
	}
}

// Stop stops measuring time and accumulates the elapsed session.
func (d *TimeCounter) Stop() {
	d.lock.Lock()
	defer d.lock.Unlock()
	if d.sessionStart != nil {
		d.continued += d.timeProvider.Now().Sub(*d.sessionStart)
		d.sessionStart = nil
	}
}

func newTimeCounterInternal(tp TimeProvider) *TimeCounter {
	tc := &TimeCounter{
		timeProvider: tp,
		lock:         sync.Mutex{},
	}
	return tc
}

// NewTimeCounter creates a new stopped counter using system time.
func NewTimeCounter() *TimeCounter {
	return newTimeCounterInternal(SystemTimeProvider())
}
