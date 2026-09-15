package pubsub

import (
	"context"
	"sync"
)

// PubSub provides a fan-out message stream with dynamic subscribers.
type PubSub[E interface{}] interface {
	// NewPublisher returns a channel used to publish events.
	NewPublisher() chan<- E
	// NewSubscriber returns a channel that receives published events until ctx is done.
	NewSubscriber(ctx context.Context) <-chan E
}

type subscription[E interface{}] struct {
	ch  chan E
	ctx context.Context
}

type pubSubImpl[E interface{}] struct {
	ctx            context.Context
	lock           sync.Mutex
	publishChannel chan E
	subscriptions  []*subscription[E]
	closed         bool
}

func (p *pubSubImpl[E]) NewPublisher() chan<- E {
	p.lock.Lock()
	defer p.lock.Unlock()
	if p.closed {
		panic("new publisher cannot be created, pubSub already closed")
	}
	return p.publishChannel
}

func (p *pubSubImpl[E]) NewSubscriber(ctx context.Context) <-chan E {
	p.lock.Lock()
	defer p.lock.Unlock()
	if p.closed {
		panic("new subscriber cannot be created, pubSub already closed")
	}
	sub := &subscription[E]{ch: make(chan E), ctx: ctx}
	p.subscriptions = append(p.subscriptions, sub)
	go func() {
		<-ctx.Done()
		p.lock.Lock()
		defer p.lock.Unlock()
		for i, s := range p.subscriptions {
			if s == sub {
				p.subscriptions = append(p.subscriptions[:i], p.subscriptions[i+1:]...)
				close(sub.ch)
				return
			}
		}
		// not found: close() has already closed the channel
	}()
	return sub.ch
}

func (p *pubSubImpl[E]) close() {
	p.lock.Lock()
	defer p.lock.Unlock()
	if !p.closed {
		p.closed = true
		close(p.publishChannel)
		for _, sub := range p.subscriptions {
			close(sub.ch)
		}
		p.subscriptions = make([]*subscription[E], 0)
	}
}

func (p *pubSubImpl[E]) push(e E) {
	p.lock.Lock()
	defer p.lock.Unlock()
	for _, sub := range p.subscriptions {
		select {
		case sub.ch <- e:
		case <-sub.ctx.Done():
			// subscriber is leaving; it will be removed once the lock is released
		case <-p.ctx.Done():
			return
		}
	}
}

// NewPubSub creates a PubSub bound to ctx.
// Canceling ctx closes publishers and subscribers.
func NewPubSub[E interface{}](ctx context.Context) PubSub[E] {
	pubCh := make(chan E)
	subs := make([]*subscription[E], 0)
	ps := &pubSubImpl[E]{
		ctx:            ctx,
		lock:           sync.Mutex{},
		publishChannel: pubCh,
		subscriptions:  subs,
	}
	go func() {
		defer func() {
			ps.close()
		}()
		for {
			select {
			case <-ctx.Done():
				return
			case e := <-pubCh:
				ps.push(e)
			}
		}
	}()
	return ps
}
