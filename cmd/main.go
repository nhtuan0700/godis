package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"
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

	fmt.Printf("Listening on host %s\n", host)

	for {
		// conn == socket == dedicated communication channel
		conn, err := lister.Accept()
		if err != nil {
			log.Fatalf("failed to accept \n")
			continue
		}

		go process(conn)
	}
}
