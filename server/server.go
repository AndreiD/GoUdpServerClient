package main

import (
	"net"
	"log"
	"fmt"
	"time"
)

const serverPort = 1337

func main() {

	server := NewServer()

	serverAddr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf(":%d", serverPort))
	errorCheck(err, "main", true)

	server.connection, err = net.ListenUDP("udp4", serverAddr)
	errorCheck(err, "main", true)

	defer server.connection.Close()

	log.Printf("Starting UDP Server, listening at %s", server.connection.LocalAddr())

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	go server.readFromSocket()
	go server.processPackets()
	go server.processMessages()

	//block indefinitely. the server will respond when someone writes to it
	select {}

}
