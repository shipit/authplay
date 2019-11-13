package main

import (
	"log"

	"github.com/shipit/authplay"
)

func main() {
	server, err := authplay.NewDefaultServer()
	if err != nil {
		log.Fatalf("%#v", err)
	}
	server.Start()
}
