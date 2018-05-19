package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/vmihailenco/msgpack"
	"net"
	"time"
	"os"
	"io"
	"path/filepath"
)

const ALIVE_CHECK_TIME = time.Second * 10

func (c *Client) setupConnection(address string) {
	addr, err := net.ResolveUDPAddr("udp4", address)

	errorCheck(err, "setupConnection", true)
	log.Printf("> server address: %s ... connecting ", addr.String())

	conn, err := net.DialUDP("udp4", nil, addr)
	c.connection = conn

	//also listen from requests from the server on a random port
	listeningAddress, err := net.ResolveUDPAddr("udp4", ":0")
	errorCheck(err, "setupConnection", true)
	log.Printf("...CONNECTED! ")

	conn, err = net.ListenUDP("udp4", listeningAddress)
	errorCheck(err, "setupConnection", true)

	log.Printf("listening on: local:%s\n", conn.LocalAddr())

}

func (c *Client) readFromSocket(buffersize int) {
	for {
		var b = make([]byte, buffersize)
		n, addr, err := c.connection.ReadFromUDP(b[0:])
		errorCheck(err, "readFromSocket", false)

		b = b[0:n]

		if n > 0 {
			pack := packet{b, addr}
			select {
			case c.packets <- pack:
				continue
			case <-c.kill:
				break
			}
		}

		//TODO: is this needed ?
		select {
		case <-c.kill:
			break
		default:
			continue
		}
	}
}

func (c *Client) processPackets() {
	for pack := range c.packets {
		var msg Message
		err := msgpack.Unmarshal(pack.bytes, &msg)
		errorCheck(err, "processPackets", false)
		c.messages <- msg

		//log.Println("<<<  CLIENT GOT")
		//spew.Dump(msg)
	}
}

func (c *Client) processMessages() {
	for msg := range c.messages {
		switch msg.Opcode {
		case 0:
			log.Printf("RRQ for file %s", msg.Filename)
		case 1:
			log.Printf("WRQ for file %s", msg.Filename)
		case 2:
			log.Printf("Data for file %s", msg.Filename)
		case 3:
			log.Printf("Acknowledgment for file %s", msg.Filename)

			//we got AKN, start sending or receiving ?

			//TODO: maybe refactor here...
			dir, _ := os.Getwd()
			fullFilePath := dir + "/aftp/client/outgoing/" + msg.Filename

			log.Println("fullFilePath:" + fullFilePath)
			if _, err := os.Stat(fullFilePath); err == nil {
				log.Info("file " + msg.Filename + " exists, sending it to the server")
				c.sendFileToServer(fullFilePath)
			} else {
				log.Info("file " + msg.Filename + " doesn't exist, waiting for the server to send it")
			}

		case 4:
			log.Printf("Error for file %s [%s]", msg.Filename, string(msg.Message))
		default:
			log.Warnln("incorrect or not implemented opcode")
		}
	}
}

func (c *Client) sendFileToServer(fullPathFile string) {

	file, err := os.Open(fullPathFile)
	errorCheck(err, "sendFileToServer", true)

	buffer := make([]byte, opts.Buffer)
	for {
		if _, err := file.Read(buffer); err == io.EOF {
			break
		}
		log.Printf("Sending > %d\n", len(string(buffer)))
		c.Send(DATA, filepath.Base(fullPathFile), buffer)
	}

}

func (c *Client) Send(opcode MessageType, filename string, payload []byte) {

	msg := Message{
		opcode, filename, payload,
	}

	//log.Println(">>> CLIENT SENDING >>> ")
	//spew.Dump(msg)

	b, err := msgpack.Marshal(msg)
	errorCheck(err, "Send", false)

	_, err = c.connection.Write(b)
	errorCheck(err, "Send", false)
}
