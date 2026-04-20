package main

import (
	"log"
	"net/http"
	_ "net/http/pprof" // for profiling
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
		// Single threaded listener server
		// err := server.RunIOMultiplexingServer(&wg)
		// if err != nil {
		// 	log.Fatal(err)
		// }

		// Multi-threaded listener server
		s, err := server.NewServer()
		if err != nil {
			log.Fatal(err)
		}
		// err = s.StartSingleListener(&wg)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		err = s.StartMultiListeners(&wg)
		if err != nil {
			log.Fatal(err)
		}
	}()

	go server.WaitForSignal(&wg, signals)

	// Expose the /debug/pprof endpoints on a separate goroutine
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	wg.Wait()
}
