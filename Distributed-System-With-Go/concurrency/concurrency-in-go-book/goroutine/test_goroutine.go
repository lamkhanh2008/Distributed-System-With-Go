package main

import (
	"fmt"
)

//	func main() {
//		// go sayHello()
//		// time.Sleep(100 * time.Millisecond)
//		// continue doing other things
//		var wg sync.WaitGroup
//		for _, salutation := range []string{"hello", "greetings", "good day"} {
//			wg.Add(1)
//			go func() {
//				defer wg.Done()
//				fmt.Println(salutation)
//			}()
//		}
//		wg.Wait()
//	}
func sayHello() {
	fmt.Println("hello")
}

func main() {
	TestOnce()
}
