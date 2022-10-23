package collector

import (
	"github.com/m-szalik/goutils"
	"sync"
	"time"
)

type TimeCounter struct {
	timeProvider goutils.TimeProvider
	sessionStart *time.Time
	continued    time.Duration
	lock         sync.Mutex
}

func (d *TimeCounter) Reset() time.Duration {
	now := d.timeProvider.Now()
	if d.sessionStart != nil {
		d.continued += now.Sub(*d.sessionStart)
		d.sessionStart = &now
	}
	dur := d.continued
	d.continued = time.Duration(0)
	return dur
}

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

func (d *TimeCounter) Start() {
	if d.sessionStart == nil {
		d.lock.Lock()
		defer d.lock.Unlock()
		now := d.timeProvider.Now()
		d.sessionStart = &now
	}
}

func (d *TimeCounter) Stop() {
	if d.sessionStart != nil {
		d.lock.Lock()
		defer d.lock.Unlock()
		d.continued += d.timeProvider.Now().Sub(*d.sessionStart)
		d.sessionStart = nil
	}
}

func newTimeCounterInternal(tp goutils.TimeProvider) *TimeCounter {
	tc := &TimeCounter{
		timeProvider: tp,
		lock:         sync.Mutex{},
	}
	return tc
}

func NewTimeCounter() *TimeCounter {
	return newTimeCounterInternal(goutils.SystemTimeProvider())
}
