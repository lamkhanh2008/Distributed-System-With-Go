package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// func ContextGo() {
// 	var Canceled = errors.New("context canceled")
// 	var DeadlineExceeded error = deadlineExceededError{}
// 	type CancelFunc
// 	type Context
// 	func Background()Context
// }

// type Context interface {
// 	Deadline() (deadline time.Time, ok bool)
// 	Done() <-chan struct{}
// 	Err() error
// 	Value(key interface{}) interface{}
// }

// func WithCancel(parent Context) (ctx Context, cancel CancelFunc)
// func WithDeadline(parent Context, deadline time.Time) (Context, CancelFunc)
// func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc)

func main() {
	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// done := make(chan interface{})
	// defer close(done)

	wg.Add(1)
	go func() {
		defer wg.Done()
		// if err := printGreeting(done); err != nil {
		if err := printGreeting(ctx); err != nil {
			fmt.Printf("cannot print greeting: %v\n", err)
			cancel()

		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		// if err := printFarewell(done); err != nil {
		if err := printFarewell(ctx); err != nil {
			fmt.Printf("cannot print farewell: %v\n", err)
			cancel()
		}
	}()
	wg.Wait()
}

// func printGreeting(done <-chan interface{}) error {
func printGreeting(ctx context.Context) error {
	greeting, err := genGreeting(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("%s world!\n", greeting)
	return nil
}
func printFarewell(ctx context.Context) error {
	farewell, err := genFarewell(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("%s world!\n", farewell)
	return nil
}

// func genGreeting(done <-chan interface{}) (string, error) {
func genGreeting(ctx context.Context) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	switch locale, err := locale(ctx); {
	case err != nil:
		return "", err
	case locale == "EN/US":
		return "hello", nil
	}
	return "", fmt.Errorf("unsupported locale")
}

// func genFarewell(done <-chan interface{}) (string, error) {
func genFarewell(ctx context.Context) (string, error) {
	switch locale, err := locale(ctx); {
	case err != nil:
		return "", err
	case locale == "EN/US":
		return "goodbye", nil
	}
	return "", fmt.Errorf("unsupported locale")
}

// func locale(done <-chan interface{}) (string, error) {
func locale(ctx context.Context) (string, error) {
	// fmt.Println("sds")
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case <-time.After(1 * time.Minute):

	}
	return "EN/US", nil
}

func locale_with_deadline(ctx context.Context) (string, error) {
	if deadline, ok := ctx.Deadline(); ok {
		if deadline.Sub(time.Now().Add(1*time.Minute)) <= 0 {
			return "", context.DeadlineExceeded
		}
	}
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case <-time.After(1 * time.Minute):
	}
	return "EN/US", nil

}
