package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	threadpool "github.com/nhtuan0700/godis/ThreadPool/theadpool"
)

func main() {
	host := ":3000"
	lister, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatalf("failed to listen on %s\n", host)
	}
	defer lister.Close()

	fmt.Printf("Listening on host %s\n", host)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// 1 pool with 2 threads
	pool := threadpool.NewPool(2)
	pool.Start()
	go func() {
		for {
			conn, err := lister.Accept()
			if err != nil {
				if pool.IsClosed() {
					log.Println("Pool is closed, shutting down")
					return
				}
				log.Printf("failed to accept: %v \n", err)
				continue
			}
			go pool.AddJob(conn)
		}
	}()

	<-sigChan
	log.Println("Received shutdown signal, shutting down gracefully...")
	pool.Close()
	lister.Close()
}
