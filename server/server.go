package main

import (
	"net"
	"time"
	log "github.com/Sirupsen/logrus"
	"github.com/jessevdk/go-flags"
)

var opts struct {
	Port   int  `short:"p" long:"port" default:":1337" description:"Port to listen on"`
	Buffer int  `short:"b" long:"buffer" default:"10240" description:"max buffer size for the socket io"`
	Quiet  bool `short:"q" long:"quiet" description:"whether to print logging info or not"`
}

func init() {
	_, err := flags.Parse(&opts)
	errorCheck(err, "init", true)
	if opts.Quiet {
		log.SetLevel(log.WarnLevel)
	}
	formatter := &log.TextFormatter{
		ForceColors : true,
		FullTimestamp: true,
		TimestampFormat: "15:04:05",

	}
	log.SetFormatter(formatter)
}

func main() {
	server := NewServer()

	uconn4, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4zero, Port: opts.Port})
	errorCheck(err, "main", true)

	server.connection = uconn4
	defer server.connection.Close()

	log.Printf("Starting UDP Server, listening at %s", server.connection.LocalAddr())

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	go server.readFromSocket(opts.Buffer)
	go server.processPackets()
	go server.processMessages()

	//block indefinitely. the server will respond when someone writes to it
	select {}

}
