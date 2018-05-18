package main

import "net"


type MessageType int

const (
	ControlMessage MessageType = iota
	TextMessage
	VoiceMessage
)

//some bytes associated with an address
type packet struct {
	bytes         []byte
	returnAddress *net.UDPAddr
}

type Client struct {
	connection *net.UDPConn
	port                int

	messages chan Message
	packets  chan packet
	kill     chan bool
}

type Message struct {
	Type        MessageType
	Message     []byte
}

//create a new client.
func NewClient() *Client {
	return &Client{
		packets:  make(chan packet),
		messages: make(chan Message),
		kill:     make(chan bool),
	}
}