package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/vmihailenco/msgpack"
	"net"
	"time"
	"os"
	"io"
	"path/filepath"
	"encoding/json"
	"github.com/davecgh/go-spew/spew"
)

const ALIVE_CHECK_TIME = time.Second * 10
var done chan bool

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

	//initialize the signal for ack on data transfers


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


	done = make(chan bool, 1)

	for msg := range c.messages {
		switch msg.Opcode {
		case RRQ:
			log.Printf("RRQ for file %s with payload %s", msg.Filename, string(msg.Message))
		case WRQ:
			log.Printf("WRQ for file %s with payload %s", msg.Filename, string(msg.Message))
		case DATA:
			log.Printf("Data for file %s", msg.Filename)
			c.WriteBytesToFile(msg.Filename, msg.Message)
		case ACK:
			if string(msg.Message) == "WRQ" {
				//sending file to server

				dir, _ := os.Getwd()
				fullFilePath := dir + "/aftp/client/myfiles/" + msg.Filename
				log.Info("sending file " + msg.Filename + " to the server. Hash: " + Sha256Sum(fullFilePath))

				if _, err := os.Stat(fullFilePath); err == nil {
					go c.sendFileToServer(fullFilePath)
				}
			} else {
				done <- true
			}

		case ERROR:
			log.Printf("Error for file %s [%s]", msg.Filename, string(msg.Message))
		case SEND_COMPLETED:
			log.Printf("SEND_COMPLETED for file %s with hash: %s", msg.Filename, string(msg.Message))
		case RECEIVED_OK:
			log.Printf("RECEIVED_OK for file %s with hash: %s", msg.Filename, string(msg.Message))
		case LIST_ALL:
			log.Printf("Received a list all request from the server. Listing....")
			var files []string
			json.Unmarshal(msg.Message, &files)
			spew.Println(files)

		default:
			log.Warnln("incorrect or not implemented opcode")
		}
	}
}

func (c *Client) sendFileToServer(fullPathFile string) {

	file, err := os.Open(fullPathFile)
	if err != nil {
		log.Warn(err)
		return
	}
	defer file.Close()
	buffer := make([]byte, opts.Buffer)
	for {
		n , err := file.Read(buffer)
		if err != nil {
			if err != io.EOF {
				log.Warn(err)
			}
			break
		}
		c.Send(DATA, filepath.Base(fullPathFile), buffer[:n])
		<-done //wait for ACK to write on channel done
	}

	c.Send(SEND_COMPLETED, filepath.Base(fullPathFile), []byte(Sha256Sum(fullPathFile)))

}

func (s *Client) WriteBytesToFile(filename string, payload []byte) {
	f, err := os.OpenFile("myfiles/"+filename, os.O_APPEND|os.O_WRONLY, 0644)
	errorCheck(err, "WriteBytesToFile", false)
	_, err = f.Write(payload)
	errorCheck(err, "WriteBytesToFile", false)
	defer f.Close()
}

func (c *Client) Send(opcode MessageType, filename string, payload []byte) {

	msg := Message{
		opcode, filename, payload,
	}

	b, err := msgpack.Marshal(msg)
	errorCheck(err, "Send", false)

	_, err = c.connection.Write(b)
	errorCheck(err, "Send", false)
}
