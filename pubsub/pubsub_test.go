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

func newTestSubscriber(parent context.Context, index int, ps PubSub[int]) *testSubscriber {
	ctx, cancel := context.WithCancel(parent)
	input := ps.NewSubscriber(ctx)
	ts := &testSubscriber{
		cancel: cancel,
		data:   make([]int, 0),
	}
	go func() {
		for {
			select {
			case d := <-input:
				fmt.Printf("subscriber %d recived data: %d\n", index, d)
				ts.data = append(ts.data, d)
			case <-ctx.Done():
				return
			}
		}
	}()
	return ts
}

func TestNewPubSubMultipleSubscribers(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()
	subscribers := make([]*testSubscriber, 0)
	ps := NewPubSub[int](ctx)
	subscribers = append(subscribers, newTestSubscriber(ctx, 0, ps))
	subscribers = append(subscribers, newTestSubscriber(ctx, 1, ps))
	publishCh := ps.NewPublisher()
	publishCh <- 0
	publishCh <- 1
	delay()
	subscribers[0].cancel()
	delay()
	publishCh <- 2
	delay()
	fmt.Println(subscribers[0].data)
	fmt.Println(subscribers[1].data)
	assert.Equal(t, []int{0, 1}, subscribers[0].data)
	assert.Equal(t, []int{0, 1, 2}, subscribers[1].data)
}

func delay() {
	time.Sleep(100 * time.Millisecond)
}
