package main

import (
	"fmt"
	"net"
	"time"
)

func process(conn net.Conn) {
	buf := make([]byte, 1000)
	conn.Read(buf)
	time.Sleep(2 * time.Second)
	conn.Write([]byte("Welcome to Godis!"))
	conn.Close()
}

func main() {
	host := ":3000"
	lister, err := net.Listen("tcp", ":3000")
	if err != nil {
		fmt.Printf("failed to listen on %s\n", host)
		return
	}

	fmt.Printf("Listening on host %s\n", host)
	
	for {
		conn, err := lister.Accept()
		if err != nil {
			fmt.Printf("failed to accpet \n")
			continue
		}

		process(conn)
	}
}
