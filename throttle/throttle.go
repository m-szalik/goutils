package throttle

// Throttler forwards events from input to output according to a throttling
// strategy.
type Throttler[E any] interface {
	// Input returns channel used to submit events.
	Input() chan<- E
	// Output returns channel that emits throttled events.
	Output() <-chan E
}

type throttler[E any] struct {
	input  chan E
	output chan E
}

func (t *throttler[E]) Input() chan<- E {
	return t.input
}

func (t *throttler[E]) Output() <-chan E {
	return t.output
}
