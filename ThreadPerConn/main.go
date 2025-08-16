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
	log.Println("Handle conn from ", conn.RemoteAddr())
	for {
		conn.SetReadDeadline(time.Now().Add(time.Minute))
		cmd, err := readCommand(conn)
		if err != nil {
			netErr, ok := err.(net.Error)
			switch {
			case ok && netErr.Timeout():
				log.Println("Read timeout")
			case err == io.EOF:
				log.Printf("client %s disconnected", conn.RemoteAddr())
			default:
				log.Printf("read error from %s: %v", conn.RemoteAddr(), err)
			}
			return
		}

		// process
		time.Sleep(time.Second)
		// Reply
		// simple http response
		cmd = fmt.Sprintf("HTTP/1.1 200 OK \r\n\r\nWelcome to Godis! command: %s \r\n", string(cmd))
		if err := respond(cmd, conn); err != nil {
			log.Println("err write: ", err)
		}
	}
}

func readCommand(conn net.Conn) (string, error) {
	var buf []byte = make([]byte, 512)
	n, err := conn.Read(buf)
	if err != nil {
		return "", err
	}

	return string(buf[:n]), nil
}

func respond(cmd string, conn net.Conn) error {
	if _, err := conn.Write([]byte(cmd)); err != nil {
		return err
	}

	return nil
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
