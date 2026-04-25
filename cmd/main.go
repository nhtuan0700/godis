package main

import (
	"log"
	"net/http"
	_ "net/http/pprof" // for profiling
	"os"
	"os/signal"
	"syscall"

	"github.com/nhtuan0700/godis/internal/server"
)

func main() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	s, err := server.NewServer()
	if err != nil {
		log.Fatal(err)
	}

	// Expose the /debug/pprof endpoints on a separate goroutine
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	if err := s.StartMultiListeners(); err != nil {
		log.Fatal(err)
	}

	s.WaitForSignal(signals)

	// wg := sync.WaitGroup{}
	// wg.Add(2)

	// go func() {
	// 	// Single threaded listener server
	// 	err := server.RunIOMultiplexingServer(&wg)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// }()

	// go server.WaitForSignal(&wg, signals)
}
