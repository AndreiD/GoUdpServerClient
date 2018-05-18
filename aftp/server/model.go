package server

import "net"

type MessageType int

const (
	ControlMessage MessageType = iota
	DataMessage
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
