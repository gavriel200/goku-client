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

const ping byte = 0xFF

type client struct {
	queueName   []byte
	queueLength uint8
}

type Consumer struct {
	client
}

type Sender struct {
	client
	conn net.Conn
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
	defer conn.Close()
	defer close(ch)

	connectionData := []byte{CONSUMER, c.client.queueLength}
	connectionData = append(connectionData, c.client.queueName...)
	conn.Write(connectionData)

	for {
		data := make([]byte, 1)
		_, err := conn.Read(data)
		if err != nil {
			fmt.Println(err, "connection closed")
			break
		}
		if data[0] == ping {
			continue
		}
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
	conn, err := net.Dial("tcp", ":8888")
	if err != nil {
		return nil, errors.New("err")
	}
	connectionData := []byte{SENDER, queueLength}
	connectionData = append(connectionData, queueName...)
	conn.Write(connectionData)
	return &Sender{client{queueName, queueLength}, conn}, nil
}

func (s *Sender) Send(data []byte) {
	s.conn.Write(data)
}
