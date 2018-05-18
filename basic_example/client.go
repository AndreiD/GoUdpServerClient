package main

import (
	"fmt"
	"net"
	"strconv"
	"time"
	"math/rand"
)

const serverAddress = "localhost:1337"

func main(){
	fmt.Println("A  Basic UDP Client Example")
	fmt.Println("Sending a random nr. to a server every 3 seconds")

	ServerAddr,err := net.ResolveUDPAddr("udp", serverAddress)
	if err != nil{
		panic(err)
	}

	//accept incoming messages from the server
	LocalAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	if err != nil{
		panic(err)
	}

	Conn, err := net.DialUDP("udp", LocalAddr, ServerAddr)
	if err != nil{
		panic(err)
	}

	defer Conn.Close()

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	ticker := time.NewTicker(3 * time.Second)
	go func() {
		for range ticker.C {
			msg := strconv.Itoa(r1.Intn(1000))
			buf := []byte(msg)

			fmt.Println("Sending > " + string(buf))
			_,err := Conn.Write(buf)
			if err != nil {
				fmt.Println(msg, err)
			}

			//after we send, we wait a response from the server
			buf = make([]byte, 1024)
			n, _, err := Conn.ReadFromUDP(buf)
			if err != nil {
				fmt.Println("Error on receiving: ",err)
			} else{
				fmt.Println("Received from server: ",string(buf[0:n]))
			}

		}
	}()

	//block forever
	select {}
}
