package pubsub

import (
	"context"
	"sync"
)

type PubSub[E interface{}] interface {
	NewPublisher() chan<- E
	NewSubscriber(ctx context.Context) <-chan E
}

type pubSubImpl[E interface{}] struct {
	lock           sync.Mutex
	publishChannel chan E
	subscriptions  []chan E
	closed         bool
}

func (p *pubSubImpl[E]) NewPublisher() chan<- E {
	if p.closed {
		panic("new publisher cannot be created, pubSub already closed")
	}
	return p.publishChannel
}

func (p *pubSubImpl[E]) NewSubscriber(ctx context.Context) <-chan E {
	if p.closed {
		panic("new subscriber cannot be created, pubSub already closed")
	}
	p.lock.Lock()
	defer p.lock.Unlock()
	ch := make(chan E)
	p.subscriptions = append(p.subscriptions, ch)
	go func() {
		<-ctx.Done()
		p.lock.Lock()
		defer p.lock.Unlock()
		for i, c := range p.subscriptions {
			if c == ch {
				p.subscriptions = append(p.subscriptions[:i], p.subscriptions[i+1:]...)
				break
			}
		}
		close(ch)
	}()
	return ch
}

func (p *pubSubImpl[E]) close() {
	p.lock.Lock()
	defer p.lock.Unlock()
	if !p.closed {
		p.closed = true
		close(p.publishChannel)
		for _, subCh := range p.subscriptions {
			close(subCh)
		}
		p.subscriptions = make([]chan E, 0)
	}
}

func (p *pubSubImpl[E]) push(e E) {
	p.lock.Lock()
	defer p.lock.Unlock()
	for _, sub := range p.subscriptions {
		sub <- e
	}
}

// NewPubSub - create new Publisher-Subscribers pair
func NewPubSub[E interface{}](ctx context.Context) PubSub[E] {
	pubCh := make(chan E)
	subs := make([]chan E, 0)
	ps := &pubSubImpl[E]{
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
