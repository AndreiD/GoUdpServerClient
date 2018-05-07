package main

import "log"

func errorCheck(err error, where string, kill bool) {
	if err != nil {
		log.Printf("Error: %s @ %s", err.Error(), where)
		if kill {
			log.Fatalln("Script Terminated")
		}
	}
}
