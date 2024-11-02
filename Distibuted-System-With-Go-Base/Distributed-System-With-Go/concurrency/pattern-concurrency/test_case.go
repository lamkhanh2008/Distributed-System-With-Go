package main

import (
	"fmt"
	"math/rand"
)

func RepeatGenerator() {
	repeat := func(done <-chan interface{}, fn func() interface{}) <-chan interface{} {
		valueStream := make(chan interface{})
		go func() {
			defer close(valueStream)
			for {
				select {
				case <-done:
					return
				case valueStream <- fn():
				}
			}
		}()
		return valueStream
	}

	take := func(done <-chan interface{}, valueStream <-chan interface{}, num int) <-chan interface{} {
		takeStream := make(chan interface{})
		go func() {
			defer close(takeStream)
			for i := 0; i < num; i++ {
				select {
				case <-done:
					return
				case takeStream <- <-valueStream:
				}
			}
		}()
		return takeStream
	}

	done := make(chan interface{})
	rand := func() interface{} {
		return rand.Int()
	}
	defer close(done)
	for num := range take(done, repeat(done, rand), 10) {
		fmt.Println(num)
	}
}

// func ToString() {
// 	toString := func(done <-chan interface{}, valueStream <-chan interface{}) <-chan interface{} {
// 		stringProcess := make(chan interface{})
// 		defer close(stringProcess)
// 		go func() {
// 			for i := range valueStream {
// 				select {
// 				case <-done:
// 					return
// 				case stringProcess <- i.(string):
// 				}
// 			}
// 		}()
// 		return stringProcess
// 	}

// }

// func main() {
// 	RepeatGenerator()
// }

// func RepeatwithFunc(){
// 	repeat := func(done <- chan interface{}, fn func() interface{})
// }
