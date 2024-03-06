package pubsub

import (
	"context"
	"fmt"
	"time"
)

func ExamplePubSub() {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()
	pubSub := NewPubSub[int](ctx)
	for i := 0; i < 4; i++ {
		// create subscriber
		go func(subscriberId int) {
			defer fmt.Printf("Subscriber %d - finised.\n", subscriberId)
			subscriberCtx, subscriberCancel := context.WithCancel(ctx)
			subscriberChannel := pubSub.NewSubscriber(subscriberCtx)
			for {
				select {
				case <-subscriberCtx.Done():
					subscriberCancel()
					return
				case number := <-subscriberChannel:
					fmt.Printf("Subscriber %d - got: number %d\n", subscriberId, number)
					if number == subscriberId {
						fmt.Printf("Subscriber %d - closes itself\n", subscriberId)
						subscriberCancel() // close itself, will not receive any more messages
					}
				}
			}
		}(i)
	}

	// create publisher
	publisherChannel := pubSub.NewPublisher()
	for msg := 0; msg < 3; msg++ {
		publisherChannel <- msg
		time.Sleep(100 * time.Millisecond)
	}
}
