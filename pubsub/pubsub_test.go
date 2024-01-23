package pubsub

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type testSubscriber struct {
	cancel context.CancelFunc
	data   []int
}

func newTestSubscriber(parent context.Context, name int, ps PubSub[int]) *testSubscriber {
	ctx, cancel := context.WithCancel(parent)
	input := ps.NewSubscriber(ctx)
	ts := &testSubscriber{
		cancel: cancel,
		data:   make([]int, 0),
	}
	go func() {
		for {
			select {
			case d, ok := <-input:
				if !ok {
					return
				}
				fmt.Printf("test-subscriber %c recived message '%d'\n", name, d)
				ts.data = append(ts.data, d)
			case <-ctx.Done():
				return
			}
		}
	}()
	return ts
}

func delay() {
	time.Sleep(100 * time.Millisecond)
}

func TestNewPubSubMultipleSubscribers(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()
	subscribers := make([]*testSubscriber, 0)
	ps := NewPubSub[int](ctx)
	subscribers = append(subscribers, newTestSubscriber(ctx, 'A', ps))
	subscribers = append(subscribers, newTestSubscriber(ctx, 'B', ps))
	publishCh := ps.NewPublisher()
	publishCh <- 0
	publishCh <- 1
	delay()
	subscribers[0].cancel()
	delay()
	publishCh <- 2
	delay()
	assert.Equal(t, []int{0, 1}, subscribers[0].data)
	assert.Equal(t, []int{0, 1, 2}, subscribers[1].data)
	subscribers = append(subscribers, newTestSubscriber(ctx, 'C', ps))
	delay()
	publishCh <- 3
	delay()
	assert.Equal(t, []int{0, 1}, subscribers[0].data)
	assert.Equal(t, []int{0, 1, 2, 3}, subscribers[1].data)
	assert.Equal(t, []int{3}, subscribers[2].data)
}

func TestNewPubSubShutdown(t *testing.T) {
	subCtx, subCancel := context.WithCancel(context.TODO())
	defer subCancel()
	subscribers := make([]*testSubscriber, 0)
	ctx, cancel := context.WithCancel(context.TODO())
	ps := NewPubSub[int](ctx)
	subscribers = append(subscribers, newTestSubscriber(subCtx, 'A', ps))
	subscribers = append(subscribers, newTestSubscriber(subCtx, 'B', ps))
	publishCh := ps.NewPublisher()
	publishCh <- 0
	publishCh <- 1
	delay()
	cancel()
	assert.Equal(t, []int{0, 1}, subscribers[0].data)
	assert.Equal(t, []int{0, 1, 2}, subscribers[1].data)
}
