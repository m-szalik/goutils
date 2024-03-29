package goutils

import (
	"fmt"
	"sync"
	"time"
)

// StopWatch - stopWatch utility
type StopWatch interface {
	// Start - start a stopWatch
	Start() StopWatch
	// Stop - stop a stop watch
	Stop()
	// Reset - reset the timer. If the timer was running it is still running after the reset.
	Reset()
	GetDuration() time.Duration
}

type stopWatch struct {
	timeProvider TimeProvider
	lock         sync.Mutex
	duration     time.Duration
	lastStart    *time.Time
}

func (s *stopWatch) Start() StopWatch {
	s.lock.Lock()
	defer s.lock.Unlock()
	now := s.timeProvider.Now()
	s.lastStart = &now
	return s
}

func (s *stopWatch) Stop() {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.lastStart != nil {
		dur := s.timeProvider.Now().Sub(*s.lastStart)
		s.lastStart = nil
		s.duration += dur
	}
}

func (s *stopWatch) Reset() {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.lastStart != nil {
		now := s.timeProvider.Now()
		s.lastStart = &now
	}
	s.duration = time.Duration(0)
}

func (s *stopWatch) GetDuration() time.Duration {
	dur, _ := s.durationCalculation()
	return dur
}

func (s *stopWatch) String() string {
	dur, ongoing := s.durationCalculation()
	return fmt.Sprintf("%s%s", dur, BoolToStr(ongoing, "+", ""))
}

func (s *stopWatch) durationCalculation() (time.Duration, bool) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.lastStart != nil {
		dur := s.timeProvider.Now().Sub(*s.lastStart)
		return s.duration + dur, true
	} else {
		return s.duration, false
	}
}

// NewStopWatch - create new StopWatch.
// New instance is not started by default.
func NewStopWatch() StopWatch {
	return NewStopWatchWithTimeProvider(SystemTimeProvider())
}

// NewStopWatchWithTimeProvider - create new StopWatch.
// New instance is not started by default.
func NewStopWatchWithTimeProvider(timeProvider TimeProvider) StopWatch {
	return &stopWatch{
		timeProvider: timeProvider,
		lock:         sync.Mutex{},
		duration:     time.Duration(0),
		lastStart:    nil,
	}
}
