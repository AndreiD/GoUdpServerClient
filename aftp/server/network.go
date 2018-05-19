package main

import (
	"net"
	log "github.com/Sirupsen/logrus"
	"github.com/vmihailenco/msgpack"
	"os"
)

func (s *Server) setupServerConnection(address string) {

	//also listen from requests from the server on a random port
	listeningAddress, err := net.ResolveUDPAddr("udp4", ":0")
	errorCheck(err, "setupConnection", true)
	log.Info("...CONNECTED! ")

	s.connection, err = net.ListenUDP("udp4", listeningAddress)
	errorCheck(err, "setupConnection", true)

	log.Printf("listening on: local:%s\n", s.connection.LocalAddr())

}

func (s *Server) readFromSocket(buffersize int) {
	for {
		var b = make([]byte, buffersize)
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

		//TODO: is this needed ?
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

		//log.Println("<<<  SERVER GOT")
		//spew.Dump(msg)
	}
}

func (s *Server) processMessages() {
	for msg := range s.messages {
		switch msg.Opcode {
		case 0:
			log.Printf("RRQ for file %s", msg.Filename)
		case 1:
			log.Printf("WRQ for file %s", msg.Filename)

			CreateDirIfNotExist("incoming")
			//will replace it if already exists
			var file, err = os.Create("incoming" + string(os.PathSeparator) + msg.Filename)
			errorCheck(err, "creating a new file", false)
			defer file.Close()

			s.Send(ACK, msg.Filename, nil)

		case 2:
			log.Printf("Data for file %s", msg.Filename)
			s.WriteBytesToFile(msg.Filename, msg.Message)
		case 3:
			log.Printf("Acknowledgment for file %s", msg.Filename)
		case 4:
			log.Printf("Error for file %s [%s]", msg.Filename, string(msg.Message))
		default:
			log.Warnln("incorrect or not implemented opcode")
		}
	}
}

func (s *Server) WriteBytesToFile(filename string, payload []byte) {
	f, err := os.OpenFile("incoming/"+ filename, os.O_APPEND|os.O_WRONLY, 0644)
	errorCheck(err, "WriteBytesToFile", false)
	_, err = f.Write(payload)
	errorCheck(err, "WriteBytesToFile", false)
	defer f.Close()
}
func (s *Server) Send(opcode MessageType, filename string, payload []byte) {

	msg := Message{
		opcode, filename, payload,
	}

	//log.Println(">>> SERVER SENDING >>> ")
	//spew.Dump(msg)

	b, err := msgpack.Marshal(msg)
	errorCheck(err, "Send", false)

	_, err = s.connection.WriteToUDP(b, s.client)
	errorCheck(err, "Send", false)
}
