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

type Message struct {
	Type    MessageType
	Message []byte
}

//"Client" Part
type Client struct {
	connection *net.UDPConn
	port       int

	messages chan Message
	packets  chan packet
	kill     chan bool
}

func NewClient() *Client {
	return &Client{
		packets:  make(chan packet),
		messages: make(chan Message),
		kill:     make(chan bool),
	}
}

//"SERVER" Part

type Server struct {
	connection *net.UDPConn
	client     *net.UDPAddr //or use map with an uuid

	messages chan Message
	packets  chan packet
	kill     chan bool
}

func NewServer() *Server {
	return &Server{
		packets:  make(chan packet),
		messages: make(chan Message),
		kill:     make(chan bool),
	}
}
