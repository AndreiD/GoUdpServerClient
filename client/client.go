package main

import (
	"time"
)

const serverAddress = "localhost:1337"

func main() {

	client := NewClient()
	client.setupConnection(serverAddress)

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	go client.readFromSocket()
	go client.processPackets()
	go client.processMessages()

	for {
		select {
		case <-ticker.C:
			client.Send("client says hello at " + time.Now().Format("15:04:05"))
		}
	}

}
