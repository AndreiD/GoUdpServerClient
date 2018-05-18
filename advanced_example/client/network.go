package main

import (
	log "github.com/Sirupsen/logrus"
	"net"

	"github.com/vmihailenco/msgpack"
	"time"
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
	}
}

func (c *Client) processMessages() {
	for msg := range c.messages {
		if msg.Type == TextMessage {
			log.Printf("Received TXT : %s", msg.Message)
		}
		if msg.Type == VoiceMessage {
			panic("todo:// voice message :)")
		}
	}
}

func (c *Client) Send(message string) {

	msg := Message{
		Type:    TextMessage,
		Message: []byte(message),
	}

	b, err := msgpack.Marshal(msg)
	errorCheck(err, "Send", false)

	_, err = c.connection.Write(b)
	errorCheck(err, "Send", false)

}
