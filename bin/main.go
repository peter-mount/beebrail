package main

import (
	"github.com/peter-mount/beebrail/server"
	"github.com/peter-mount/beebrail/server/status"
	"github.com/peter-mount/golib/kernel"
	"log"
)

func main() {
	err := kernel.Launch(
		&server.Server{},
		&status.Status{},
		&server.Telnet{},
	)
	if err != nil {
		log.Fatal(err)
	}
}
