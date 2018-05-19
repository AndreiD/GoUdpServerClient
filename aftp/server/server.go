package main

import (
	"net"
	log "github.com/Sirupsen/logrus"
	"github.com/jessevdk/go-flags"
)

var opts struct {
	Port  int  `short:"p" long:"port" default:"6969" description:"port to listen to."`
	Buffer  int  `short:"b" long:"buffer" default:"1024" description:"buffer size. default 1024"`
	Quiet bool `short:"q" long:"quiet" description:"print less logging information"`
}

func init() {
	_, err := flags.Parse(&opts)
	errorCheck(err, "init", true)
	if opts.Quiet {
		log.SetLevel(log.WarnLevel)
	}
	formatter := &log.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "15:04:05",
	}
	log.SetFormatter(formatter)
}

func main() {
	server := NewServer()

	udpConn, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4zero, Port: opts.Port})
	errorCheck(err, "main", true)

	server.connection = udpConn
	defer server.connection.Close()

	log.Printf("Starting AFTP Server %s:%d", GetLocalIP(), opts.Port)



	go server.readFromSocket(opts.Buffer)
	go server.processPackets()
	go server.processMessages()

	//block indefinitely. the server will respond when someone writes to it
	select {}

}
