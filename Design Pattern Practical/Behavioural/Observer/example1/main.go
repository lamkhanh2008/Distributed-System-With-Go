package main

import (
	"observer/example1/observer"
	"observer/example1/subject"
)

func main() {

	shirtItem := subject.NewProducer("ADIDAS", true)

	observerFirst := &observer.Customer{Id: "10"}
	observerSecond := &observer.Customer{Id: "11"}
	shirtItem.Register(observerFirst)
	shirtItem.Register(observerSecond)
	shirtItem.NotifyAll()

}
