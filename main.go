package main

import (
	"flag"
	"log"
	"time"
)

var (
	host = flag.String("host", "localhost", "Host")
	port = flag.Int("port", 8099, "Port")
)

func init() {
	flag.Parse()
}

func main() {
	protocol := NewProtocol()
	server := NewServer(protocol)

	protocol.RegisterCommand(NewPublishCommand(server))
	protocol.RegisterCommand(NewConsumeCommand(server))

	err := server.Run(*host, *port)

	if err != nil {
		log.Fatalf("[%v] (!) %s", time.Now(), err.Error())
	}
}
