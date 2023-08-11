package goutils

import (
	"fmt"
	"sync"
	"time"
)

type StopWatch interface {
	Start() StopWatch
	Stop()
	Reset()
	GetDuration() time.Duration
}

type stopWatch struct {
	lock      sync.Mutex
	duration  time.Duration
	lastStart *time.Time
}

func (s *stopWatch) Start() StopWatch {
	s.lock.Lock()
	defer s.lock.Unlock()
	now := time.Now()
	s.lastStart = &now
	return s
}

func (s *stopWatch) Stop() {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.lastStart != nil {
		dur := time.Now().Sub(*s.lastStart)
		s.lastStart = nil
		s.duration += dur
	}
}

func (s *stopWatch) Reset() {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.lastStart = nil
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
		dur := time.Now().Sub(*s.lastStart)
		return s.duration + dur, true
	} else {
		return s.duration, false
	}
}

func NewStopWatch() StopWatch {
	return &stopWatch{
		lock:      sync.Mutex{},
		duration:  time.Duration(0),
		lastStart: nil,
	}
}
