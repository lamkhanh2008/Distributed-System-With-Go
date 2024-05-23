package main

import "fmt"

func OrDone(done <-chan interface{}, stream <-chan interface{}) chan interface{}

func BridgeChannel() {
	done := make(chan interface{})
	chanStream := chanStream()
	bridge := func(done <-chan interface{}, chanStream <-chan <-chan interface{}) <-chan interface{} {
		valStream := make(chan interface{})
		go func() {
			defer close(valStream)
			for {
				var stream <-chan interface{}
				select {
				case maybeStream, ok := <-chanStream:
					if ok == false {
						return
					}
					stream = maybeStream

				case <-done:
					return
				}

				for val := range OrDone(done, stream) {
					select {
					case valStream <- val:
					case <-done:
					}
				}

			}

		}()
		return valStream
	}

	for i := range bridge(done, chanStream) {
		fmt.Println(i)
	}
}

func chanStream() <-chan <-chan interface{} {

	chanStream := make(chan (<-chan interface{}))
	go func() {
		defer close(chanStream)
		for i := 0; i < 10; i++ {
			stream := make(chan interface{}, 1)
			stream <- 1
			close(stream)
			chanStream <- stream
		}
	}()
	return chanStream
}

//create series of 10 channels, each with one elemnent write to them, passes these channels into bridge function
