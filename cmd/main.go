package main

import (
	"log"

	"github.com/nhtuan0700/godis/internal/server"
)

func main() {
	err := server.RunIOMultiplexingServer()
	if err != nil {
		log.Fatal(err)
	}
}
