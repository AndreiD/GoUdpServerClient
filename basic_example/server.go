package main

import (
	"fmt"
	"net"
)

const serverPort = ":1337"

func main() {
	fmt.Println("A  Basic UDP Server Example")

	ServerAddr, err := net.ResolveUDPAddr("udp", serverPort)
	if err != nil {
		panic(err)
	}

	ServerConn, err := net.ListenUDP("udp", ServerAddr)
	if err != nil {
		panic(err)
	}
	defer ServerConn.Close()

	buf := make([]byte, 1024)
	go func() {
		for {
			n, addr, err := ServerConn.ReadFromUDP(buf)
			fmt.Println("Received ", string(buf[0:n]), " from ", addr)

			if err != nil {
				fmt.Println("Error: ", err)
			}

			//after we got something, respond with an "OK" to the client
			buf = []byte("OK")
			ServerConn.WriteToUDP(buf, addr)
		}
	}()

	fmt.Println("Waiting for clients to connect. Server port " + serverPort)

	//blocking forever
	select {}
}
