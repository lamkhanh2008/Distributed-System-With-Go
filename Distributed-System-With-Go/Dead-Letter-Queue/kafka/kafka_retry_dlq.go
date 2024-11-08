package main

import (
	"context"
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/segmentio/kafka-go"
)

/*
ProcessRetryHandler: An interface that must be implemented by types that will handle the processing of messages. It has the following methods:
Process(context.Context, kafka.Message) error: processes a message, returns an error if the processing failed.
MoveToDLQ(context.Context, kafka.Message): moves a message to the dead letter queue.
MaxRetries() int: returns the maximum number of retries for a message.
Backoff() backoff.BackOff: returns a backoff strategy for retrying a message.
*/
type ProcessRetryHandler interface {
	Process(context.Context, kafka.Message) error
	MoveToDLQ(context.Context, kafka.Message)
}

type ConsumerWithRetryOptions struct {
	Handler    ProcessRetryHandler
	Reader     *kafka.Reader
	MaxRetries int
	RetryQueue chan kafka.Message
	Backoff    backoff.BackOff
}

func NewConsumerWithRetry(ctx context.Context, options *ConsumerWithRetryOptions) {
	go func() {
		for {
			select {
			case msg := <-options.RetryQueue:
				retries := 1
				for {
					fmt.Printf("Retry %v message %v\n", retries, msg.Key)
					if retries >= options.MaxRetries {
						options.Handler.MoveToDLQ(ctx, msg)
						break
					}

					if err := options.Handler.Process(ctx, msg); err != nil {
						fmt.Printf("Error processing message, retrying: %v\n", err)
						time.Sleep(options.Backoff.NextBackOff())
						retries++
						continue
					}

					break
				}
			}
		}
	}()

	for {
		msg, err := options.Reader.ReadMessage(ctx)
		if err != nil {
			fmt.Println("Error reading message: ", err)
			return
		}

		if err := options.Handler.Process(ctx, msg); err != nil {
			options.RetryQueue <- msg
		}
	}
}
