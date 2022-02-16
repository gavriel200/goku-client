package main

import (
	"fmt"

	"github.com/gavriel200/goku-client"
)

func main() {
	ch := make(chan []byte)
	consumer, err := goku.NewConsumer("test")
	if err != nil {
		fmt.Println(err)
		return
	}
	go consumer.Listen(ch)
	for data := range ch {
		fmt.Println(data)
	}
}
