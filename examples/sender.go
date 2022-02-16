package main

import (
	"fmt"

	"github.com/gavriel200/goku-client"
)

func main() {
	sender, err := goku.NewSender("test")
	if err != nil {
		fmt.Println(err)
		return
	}
	sender.Send([]byte{10})
}
