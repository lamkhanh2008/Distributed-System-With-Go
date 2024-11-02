package main

import (
	"fmt"
	"time"
)

type Message struct {
	topic   string
	content interface{}
}

type MessageChannel chan Message

func main() {
	maxMessage := 10000
	topic := "update-user"

	messageQueue := make(chan Message, maxMessage)
	mapTopicMessage := make(map[string][]MessageChannel) // map[topic][]MessageChannel
	go run(messageQueue, mapTopicMessage)
	publish(messageQueue, topic, "user-name is update to Hung")
}

func run(messageQueue chan Message, mapTopicMessage map[string][]MessageChannel) {
	for {
		message := <-messageQueue
		listMessageChannel, ok := mapTopicMessage[message.topic]
		if ok {
			for _, messageChannel := range listMessageChannel {
				messageChannel <- message
			}
		}
	}
}

func publish(messageQueue chan Message, topic string, content string) {
	message := Message{
		topic:   topic,
		content: content,
	}
	messageQueue <- message
	fmt.Printf("%v: publish new message with topic: '%v' - content: '%v' \n", time.Now().Format("15:04:05"), message.topic, message.content)
}

func registerSubscription(mapTopicMessage map[string][]MessageChannel, topic string) MessageChannel {
	newMessageChannel := make(MessageChannel)
	value, ok := mapTopicMessage[topic]
	if ok {
		value = append(value, newMessageChannel)
		mapTopicMessage[topic] = value
	} else {
		mapTopicMessage[topic] = []MessageChannel{newMessageChannel}
	}

	return newMessageChannel
}
