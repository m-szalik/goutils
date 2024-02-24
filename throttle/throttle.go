package throttle

type Throttler[E any] interface {
	Input() chan<- E
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
