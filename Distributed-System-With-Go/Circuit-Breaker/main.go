package main

import (
	circuitbreaker "circuit_breaker/base_circuit_breaker"
	"errors"
	"fmt"
	"log"
	"math/rand/v2"
	"sync"
	"time"
)

func main() {
	cbOpts := circuitbreaker.ExtraOptions{
		Policy:              circuitbreaker.MaxFails,
		MaxFails:            circuitbreaker.ToPointer(uint64(5)),
		MaxConsecutiveFails: circuitbreaker.ToPointer(uint64(5)),
		OpenInterval:        circuitbreaker.ToPointer(50 * time.Millisecond),
	}
	cb := circuitbreaker.New(cbOpts)
	wg := &sync.WaitGroup{}
	for i := 1; i < 30; i += 1 {
		wg.Add(1)
		go makeServiceCall(i, cb, wg)
		// time.Sleep(10 * time.Millisecond)
	}

	log.Println("sent all the requests")
	wg.Wait()
	log.Println("got all the responses, exiting.")

}
func serviceMethod(id int) (string, error) {
	if val := rand.Float64(); val <= 0.5 {
		return "", errors.New("failed")
	}
	return fmt.Sprintf("[id: %d] done.", id), nil
}

func makeServiceCall(id int, cb circuitbreaker.CircuitBreaker, wg *sync.WaitGroup) {
	defer wg.Done()

	resp, err := cb.Execute(func() (interface{}, error) {
		return serviceMethod(id)
	})
	if err != nil {
		log.Printf("[id %d] got err: %s", id, err.Error())
	} else {
		log.Printf("[id %d] success: %s", id, resp)
	}
}
