package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nhtuan0700/godis/internal/threadpool"
)

func process(conn net.Conn) {
	defer conn.Close()
	// Read data from client
	buf := make([]byte, 1000)
	for {
		conn.SetReadDeadline(time.Now().Add(10 * time.Second))
		_, err := conn.Read(buf)
		if err != nil {
			netErr, ok := err.(net.Error)
			switch {
			case ok && netErr.Timeout():
				log.Println("Read timeout")
			case err == io.EOF:
				log.Printf("client %s closed connection", conn.RemoteAddr())
			default:
				log.Printf("read error from %s: %v", conn.RemoteAddr(), err)
			}
			return
		}

		// process
		time.Sleep(time.Second)
		log.Printf("Request from %s\n", conn.RemoteAddr())
		// Reply
		conn.Write([]byte("HTTP/1.1 200 OK \r\n\r\nWelcome to Godis!\r\n"))
	}
}

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
