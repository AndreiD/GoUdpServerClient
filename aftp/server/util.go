package main

import (
	log "github.com/Sirupsen/logrus"
	"net"
	"os"
	"crypto/sha256"
	"io"
	"encoding/hex"
	"path/filepath"
	"encoding/json"
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

// GetLocalIP returns the non loopback local IP of the host
func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

func CreateDirIfNotExist(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			panic(err)
		}
	}
}

// Calculate Sha256 hash
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

// WalkFunc for ListAllFiles
func visit(files *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Warn(err)
			return nil
		}
		*files = append(*files, path)
		return nil
	}
}

// List all files from a directory as marshaled json array
func ListAllFiles(root string) []byte {
	var files []string

	err := filepath.Walk(root, visit(&files))
	if err != nil {
		log.Warn(err)
		return nil
	}

	data, err := json.Marshal(files)

	if err != nil {
		log.Warn(err)
		return nil
	}
	return data
}
