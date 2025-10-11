package server

import (
	"log"
	"os"
	"sync"
	"sync/atomic"

	"github.com/nhtuan0700/godis/internal/constant"
)

func WaitForSignal(wg *sync.WaitGroup, signals chan os.Signal) {
	defer wg.Done()
	// Wait for signal in channel, it not available then wait
	<-signals
	// log.Println("Shutting down gracefully...")
	// os.Exit(0)

	// Busy loop
	log.Println("Shutting down gracefully...")
	for {
		if atomic.CompareAndSwapInt32(&serverStatus, constant.ServerStatusIdle, constant.ServerStatusShuttingDown) {
			os.Exit(0) // shutdown
		}
	}
}
