package main

import (
	"time"
	"github.com/jessevdk/go-flags"
	log "github.com/Sirupsen/logrus"
)

var opts struct {
	ServerAddr   string  `short:"s" long:"port" default:":6969" description:"server address"`
	Filename string  `short:"f" long:"file" default:"" description:"name of the file you want to send or get"`
	Type string  `short:"t" long:"type" default:"send" description:"use send or get"`
	Quiet  bool `short:"v" long:"verbose" description:"print less logging information"`
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
