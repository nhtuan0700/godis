package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"syscall"

	"github.com/nhtuan0700/godis/internal/core"
	"github.com/nhtuan0700/godis/internal/core/io_multiplexer"
)

type IOHandler struct {
	id            int
	ioMultiplexer io_multiplexer.IOMultiplexer
	server        *Server
	mu            sync.Mutex
	// Garbage Collector may close the connection if we don't keep a reference to it in the I/O handler
	// We use a map to store active connections, the key is the file descriptor of the connection
	// when running benchmark, the number of connections can be very large, the gc run quickly and close the connection before the I/O handler can read from it
	// which causes "bad file descriptor" error -> benchmark fails
	conns map[int]net.Conn
}

func NewIOHandler(id int, server *Server) (*IOHandler, error) {
	ioMultiplexer, err := io_multiplexer.CreateIOMultiplexer()
	if err != nil {
		return nil, err
	}

	ioHandler := &IOHandler{
		id:            id,
		ioMultiplexer: ioMultiplexer,
		server:        server,
		conns:         make(map[int]net.Conn),
	}

	return ioHandler, nil
}

func (h *IOHandler) AddConn(conn net.Conn) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	tcpConn, ok := conn.(*net.TCPConn)
	if !ok {
		return fmt.Errorf("failed to cast net.Conn to *net.TCPConn")
	}
	rawConn, err := tcpConn.SyscallConn()
	if err != nil {
		return err
	}

	// Add the connection's file descriptor to the I/O multiplexer for monitoring
	err = rawConn.Control(func(fd uintptr) {
		connFd := int(fd)
		log.Printf("I/O Handler %d is monitoring fd %d", h.id, connFd)
		h.conns[connFd] = conn
		h.ioMultiplexer.Monitor(io_multiplexer.Event{
			Fd: connFd,
			Op: io_multiplexer.OpRead,
		})
	})

	return err
}

// Run starts the event loop for the I/O handler
// waiting for events on monitored file descriptors and processing them
func (h *IOHandler) Run() {
	log.Printf("I/O Handler %d started\n", h.id)

	for {
		// wait for data from any fd in the monitoring list
		events, err := h.ioMultiplexer.Wait()
		if err != nil {
			log.Printf("I/O Handler %d error while waiting for events: %v\n", h.id, err)
			continue
		}

		for _, event := range events {
			connFd := event.Fd
			// log.Printf("I/O Handler %d received event on fd %d\n", h.id, connFd)

			cmd, err := readCommand(connFd)
			if err != nil {
				if err == io.EOF || err == syscall.ECONNRESET {
					log.Printf("I/O Handler %d: connection closed on fd %d\n", h.id, connFd)
				} else {
					log.Printf("Read error on fd %d: %v\n", connFd, err)
				}
				// if err := syscall.Close(connFd); err != nil {
				// 	log.Printf("Error closing fd %d: %v\n", connFd, err)
				// }
				continue
			}

			replyChan := make(chan []byte, 1)
			task := &core.Task{
				Command:   cmd,
				ReplyChan: replyChan,
			}

			// dispatch the command to the corresponding worker
			h.server.dispatch(task)

			res := <-replyChan
			if err := respond(res, connFd); err != nil {
				log.Printf("Write error on fd %d: %v\n", connFd, err)
			}
		}
	}
}
