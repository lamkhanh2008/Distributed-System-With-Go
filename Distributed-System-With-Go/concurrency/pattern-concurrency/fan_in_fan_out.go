package main

import (
	"fmt"
	"sync"
)

func FanInBasic() <-chan interface{} {
	done := make(chan interface{})
	defer close(done)
	channelResult := make(chan interface{})
	valueStreams := make([]chan interface{}, 5)
	var wg sync.WaitGroup
	fanIn := func(done <-chan interface{}, valueStream <-chan interface{}) {
		defer wg.Done()
		for i := range valueStream {
			select {
			case <-done:
				return
			case channelResult <- i:
			}
		}
	}

	wg.Add(len(valueStreams))
	for _, c := range valueStreams {

		go fanIn(done, c)

	}

	go func() {
		wg.Wait()
		close(channelResult)
	}()
	return channelResult
}

func OrDoneChannel() {
	done := make(chan interface{})
	close(done)
	valueStream := make(chan interface{})
	orDone := func(done <-chan interface{}, valueStream chan interface{}) chan string {
		resultChannel := make(chan string)
		defer close(resultChannel)
		go func() {
			for {
				select {
				case <-done:
					return
				case v, ok := <-valueStream:
					if ok == false {
						return
					}

					select {
					case valueStream <- v:
					case <-done:
					}
				}
			}
		}()
		return resultChannel
	}

	for i := range orDone(done, valueStream) {
		fmt.Println(i)
	}
}
