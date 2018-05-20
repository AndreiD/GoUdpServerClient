package main

import "net"

type MessageType int

//list of operation codes
const (
	RRQ   MessageType = iota
	WRQ
	DATA
	ACK
	ERROR
	SEND_COMPLETED
	RECEIVED_OK
	LIST_ALL
)

//some bytes associated with an address
type packet struct {
	bytes         []byte
	returnAddress *net.UDPAddr
}

type Message struct {
	Opcode   MessageType
	Filename string
	Message  []byte
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
