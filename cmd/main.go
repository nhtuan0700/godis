package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/nhtuan0700/godis/internal/server"
)

func main() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		err := server.RunIOMultiplexingServer(&wg)
		if err != nil {
			log.Fatal(err)
		}
	}()

	go server.WaitForSignal(&wg, signals)
	wg.Wait()
}
