package main

import (
	"net"
	"log"
	"fmt"
	"time"
)

const port = 1337

type Server struct {
	connection *net.UDPConn
	messages   chan string
	client     *net.UDPAddr //or use map with an uuid
}

var buffer = make([]byte, 1024)

func errorCheck(err error, where string, kill bool) {
	if err != nil {
		log.Printf("‡‡‡ Error ‡‡‡ %s | %s", err.Error(), where)
		if kill {
			log.Fatalln("Script Terminated")
		}
	}
}

func (s *Server) handleMessage() {
	var buf [512]byte

	n, addr, err := s.connection.ReadFromUDP(buf[0:])
	errorCheck(err, "handleMessage", false)

	got := string(buf[0:n])
	log.Printf("a client [%s] sent [%s]\n", addr, got)

	// respond with something ?
	s.client = addr

	s.messages <- "server says hello at " + time.Now().Format("15:04:05")
}

func (s *Server) sendMessage() {
	for {
		msg := <-s.messages
		err := s.connection.SetWriteDeadline(time.Now().Add(1 * time.Second))
		if err != nil {
			log.Printf("Seems like this client %s is gone.\n", s.client)

		}
		_, err = s.connection.WriteTo([]byte(msg), s.client)
		errorCheck(err, "sendMessage", false)
	}
}

func main() {

	var s Server
	s.messages = make(chan string, 20)

	serverAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", port))
	errorCheck(err, "main", true)

	s.connection, err = net.ListenUDP("udp", serverAddr)
	errorCheck(err, "main", true)

	defer s.connection.Close()

	log.Printf("Starting UDP Server, listening at %s", s.connection.LocalAddr())

	go s.sendMessage()

	for {
		s.handleMessage()
	}

}
