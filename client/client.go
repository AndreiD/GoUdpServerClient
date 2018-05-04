package main

import (
	"fmt"
	"net"
	"time"
	"log"
	"strings"
)

const serverAddress = "localhost"
const port = 1337

type Client struct {
	connection          *net.UDPConn
	alive               bool
	sendingMessageQueue chan string
	receiveMessages     chan string
}

func errorCheck(err error, where string, kill bool) {
	if err != nil {
		log.Printf("Error: %s @ %s", err.Error(), where)
		if kill {
			log.Fatalln("Script Terminated")
		}
	}
}

// while the client is alive, if there's a message on queue, send it
func (c *Client) sendMessage() {
	for c.alive {
		msg := <-c.sendingMessageQueue
		_, err := c.connection.Write([]byte(msg))
		errorCheck(err, "sendMessage", false)
	}

}

// handles the receiving part
func (c *Client) receiveMessage() {
	var buffer = make([]byte, 1024)
	for c.alive {
		n, err := c.connection.Read(buffer[0:])
		errorCheck(err, "receiveMessage", false)
		c.receiveMessages <- string(buffer[0:n])
	}
}

// separates the receiving from the processing
func (c *Client) processMessage() {
	for c.alive {
		msg := <-c.receiveMessages
		log.Printf("server [%s] sent [%s]\n", c.connection.RemoteAddr().String(), msg)

		//additional processing here
		if strings.HasPrefix(msg, ":q") || strings.HasPrefix(msg, ":quit") {
			fmt.Printf("%s is leaving\n", c.connection.RemoteAddr().String())
		}
	}
}

func main() {
	serverAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", serverAddress, port))
	errorCheck(err, "main", true)

	var c Client
	c.alive = true
	c.sendingMessageQueue = make(chan string)
	c.receiveMessages = make(chan string)

	c.connection, err = net.DialUDP("udp", nil, serverAddr) //localAddr is automatically chosen
	errorCheck(err, "main", true)

	defer c.connection.Close()

	log.Printf("Starting UDP Client, connected to %s (localAddress %s)", c.connection.RemoteAddr(), c.connection.LocalAddr())

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	//fire it up!
	go c.receiveMessage()
	go c.processMessage()

	go c.sendMessage()

	for {
		select {
		case <-ticker.C:
			c.sendingMessageQueue <- "client says hello at " + time.Now().Format("15:04:05")
		}
	}

}
