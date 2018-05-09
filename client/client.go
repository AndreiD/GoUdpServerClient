package main

import (
	"time"
	"github.com/jessevdk/go-flags"
	log "github.com/Sirupsen/logrus"
)

var opts struct {
	ServerAddress string `short:"s" long:"serveraddress" default:"localhost:1337" description:"The Server's Address (ex: localhost:1337)"`
	Buffer        int    `short:"b" long:"buffer" default:"10240" description:"max buffer size for the socket io"`
	Quiet         bool   `short:"q" long:"quiet" description:"whether to print logging info or not"`
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

	client := NewClient()
	client.setupConnection(opts.ServerAddress)

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	go client.readFromSocket(opts.Buffer)
	go client.processPackets()
	go client.processMessages()

	for {
		select {
		case <-ticker.C:
			client.Send("client says hello at " + time.Now().Format("15:04:05"))
		}
	}

}
