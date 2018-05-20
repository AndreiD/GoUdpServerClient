package main

import (
	log "github.com/Sirupsen/logrus"
	"os"
	"crypto/sha256"
	"io"
	"encoding/hex"
)

func errorCheck(err error, where string, kill bool) {
	if err != nil {
		if kill {
			log.WithError(err).Fatalln("Script Terminated")
		} else {
			log.WithError(err).Warnf("@ %s\n", where)
		}
	}
}

//calculate Sha256 hash
func Sha256Sum(filepath string) string {
	f, err := os.Open(filepath)
	if err != nil {
		log.Warn(err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err = io.Copy(h, f); err != nil {
		log.Warn(err)
	}

	return hex.EncodeToString(h.Sum(nil))
}
