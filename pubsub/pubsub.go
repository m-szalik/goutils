package pubsub

import (
	"context"
	"sync"
)

type PubSub[E interface{}] interface {
	Publisher() chan<- E
	NewSubscriber() <-chan E
}

type pubSubImpl[E interface{}] struct {
	lock           sync.Mutex
	publishChannel chan E
	subscriptions  []chan E
}

func (p *pubSubImpl[E]) Publisher() chan<- E {
	return p.publishChannel
}

func (p *pubSubImpl[E]) NewSubscriber() <-chan E {
	p.lock.Lock()
	defer p.lock.Unlock()
	ch := make(chan E)
	p.subscriptions = append(p.subscriptions, ch)
	return ch
}

func (p *pubSubImpl[E]) closeSubs() {
	p.lock.Lock()
	defer p.lock.Unlock()
	for _, sub := range p.subscriptions {
		close(sub)
	}
}

func (p *pubSubImpl[E]) push(e E) {
	p.lock.Lock()
	defer p.lock.Unlock()
	for _, sub := range p.subscriptions {
		sub <- e
	}
}

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
			close(pubCh)
			ps.closeSubs()
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
