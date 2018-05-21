package main

import (
	"github.com/jessevdk/go-flags"
	log "github.com/Sirupsen/logrus"
	"os"
)

var opts struct {
	ServerAddr string `short:"s" long:"port" default:"localhost:6969" description:"server address"`
	Buffer     int    `short:"b" long:"buffer" default:"10240" description:"buffer size. default 1024"`
	Filename   string `short:"f" long:"file" default:"" description:"name of the file you want to send or get"`
	Type       string `short:"t" long:"type" default:"send" description:"use send or get"`
	Quiet      bool   `short:"v" long:"verbose" description:"print less logging information"`
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

	client := NewClient()
	client.setupConnection(opts.ServerAddr)

	go client.readFromSocket(opts.Buffer)
	go client.processPackets()
	go client.processMessages()

	//sending 1mb.bin
	filetosend := "5kb.bin"
	dir, _ := os.Getwd()
	fullFilePath := dir + "/aftp/client/myfiles/" + filetosend

	client.Send(WRQ, filetosend, []byte(Sha256Sum(fullFilePath)))

	//block forever
	select {}

}
