package main

import (
	"log"

	"github.com/shipit/authplay"
)

func main() {
	// log filename and line numbers
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	server, err := authplay.NewServer()
	if err != nil {
		log.Fatalf("%#v", err)
	}
	server.Start()
}
