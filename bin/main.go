package main

import (
	"github.com/peter-mount/beebrail/server"
	"github.com/peter-mount/beebrail/server/connector"
	"github.com/peter-mount/beebrail/server/status"
	"github.com/peter-mount/golib/kernel"
	"log"
)

func main() {
	err := kernel.Launch(
		&server.Server{},
		&status.Status{},
		&connector.SerialPort{},
		&connector.Telnet{},
	)
	if err != nil {
		log.Fatal(err)
	}
}
