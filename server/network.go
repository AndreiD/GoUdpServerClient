package main

import (
	"log"
	"net"

	"github.com/vmihailenco/msgpack"
)

func (s *Server) setupServerConnection(address string) {

	//also listen from requests from the server on a random port
	listeningAddress, err := net.ResolveUDPAddr("udp4", ":0")
	errorCheck(err, "setupConnection", true)
	log.Printf("...CONNECTED! ")

	s.connection, err = net.ListenUDP("udp4", listeningAddress)
	errorCheck(err, "setupConnection", true)

	log.Printf("listening on: local:%s\n", s.connection.LocalAddr())

}

func (s *Server) readFromSocket() {
	for {
		var b = make([]byte, 1024)
		n, addr, err := s.connection.ReadFromUDP(b[0:])
		errorCheck(err, "readFromSocket", false)

		s.client = addr

		b = b[0:n]
		if n > 0 {
			pack := packet{b, addr}
			select {
			case s.packets <- pack:
				continue
			case <-s.kill:
				break
			}
		}

		select {
		case <-s.kill:
			break
		default:
			continue
		}
	}
}

func (s *Server) processPackets() {
	for pack := range s.packets {
		var msg Message
		err := msgpack.Unmarshal(pack.bytes, &msg)
		errorCheck(err, "processPackets", false)
		s.messages <- msg
	}
}

func (s *Server) processMessages() {
	for msg := range s.messages {
		if msg.Type == TextMessage {
			log.Printf("Received TXT : %s\n", msg.Message)

			//respond something
			s.Send("the server got your message")
		}
		if msg.Type == VoiceMessage {
			panic("todo:// voice message :)")
		}
	}
}

func (s *Server) Send(message string) {

	msg := Message{
		Type:    TextMessage,
		Message: []byte(message),
	}

	b, err := msgpack.Marshal(msg)
	errorCheck(err, "Send", false)

	_, err = s.connection.WriteToUDP(b, s.client)
	errorCheck(err, "Send", false)

}
