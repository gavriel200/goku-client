package goku

import (
	"errors"
	"fmt"
	"net"
)

// default configuration
const (
	DefaultURL  = "goku://localhost"
	DefaultPort = ":8888"
)

const (
	CONSUMER byte = iota
	SENDER
)

type client struct {
	queueName   []byte
	queueLength uint8
}

type Consumer struct {
	client
}

type Sender struct {
	client
}

func NewConsumer(queueNameString string) (*Consumer, error) {
	queueName := []byte(queueNameString)
	queueLengthInt := len(queueName)
	if queueLengthInt > 255 {
		return nil, errors.New("queue name cant be longer then 255 characters")
	}
	queueLength := uint8(queueLengthInt)
	return &Consumer{client{queueName, queueLength}}, nil
}

func (c *Consumer) Listen(ch chan []byte) {
	conn, err := net.Dial("tcp", ":8888")
	if err != nil {
		fmt.Println(err)
		return
	}
	connectionData := []byte{CONSUMER, c.client.queueLength}
	connectionData = append(connectionData, c.client.queueName...)
	conn.Write(connectionData)

	for {
		data := make([]byte, 1)
		conn.Read(data)
		ch <- data
	}
}

func NewSender(queueNameString string) (*Sender, error) {
	queueName := []byte(queueNameString)
	queueLengthInt := len(queueName)
	if queueLengthInt > 255 {
		return nil, errors.New("queue name cant be longer then 255 characters")
	}
	queueLength := uint8(queueLengthInt)
	return &Sender{client{queueName, queueLength}}, nil
}

func (c *Sender) Send(data []byte) {
	conn, err := net.Dial("tcp", ":8888")
	if err != nil {
		fmt.Println(err)
		return
	}
	connectionData := []byte{CONSUMER, c.client.queueLength}
	connectionData = append(connectionData, c.client.queueName...)
	conn.Write(connectionData)
	conn.Write(data)
}
