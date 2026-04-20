package server

import (
	"context"
	"log"
	"net"
	"sync"
	"syscall"

	"github.com/nhtuan0700/godis/internal/config"
)

func (s *Server) StartSingleListener(wg *sync.WaitGroup) error {
	defer wg.Done()

	// Start all I/O handlers event loops
	for _, handler := range s.ioHandlers {
		go handler.Run()
	}

	// Setup listener socket

	listener, err := net.Listen(config.Protocol, config.Address)
	if err != nil {
		return err
	}
	defer listener.Close()

	log.Printf("Server listening on %s", config.Address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		handler := s.ioHandlers[s.nextIOHandler]
		s.nextIOHandler = (s.nextIOHandler + 1) % len(s.ioHandlers)

		if err := handler.AddConn(conn); err != nil {
			log.Printf("Failed to add connection to I/O handler %d: %v", handler.id, err)
			// If adding fails, close the connection to avoid resource leak
			conn.Close()
		}
	}
}

func createReusablePortListener(network, addr string) (net.Listener, error) {
	// Create a socket with SO_REUSEPORT option
	lc := net.ListenConfig{
		Control: func(network, address string, c syscall.RawConn) error {
			var err error
			err = c.Control(func(fd uintptr) {
				err = syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_REUSEPORT, 1)
			})
			return err
		},
	}

	return lc.Listen(context.Background(), network, addr)
}

func (s *Server) StartMultiListeners(wg *sync.WaitGroup) error {
	defer wg.Done()

	// Start all I/O handlers event loops
	for _, handler := range s.ioHandlers {
		go handler.Run()
	}

	// Start a listener for each I/O handler
	for i := 0; i < config.ListenerNumber; i++ {
		go func() {
			listener, err := createReusablePortListener(config.Protocol, config.Address)
			if err != nil {
				log.Fatal(err)
			}
			defer listener.Close()
			log.Printf("Listener %d started on %s", i, config.Address)

			for {
				conn, err := listener.Accept()
				if err != nil {
					log.Printf("Failed to accept connection: %v", err)
					continue
				}

				handler := s.ioHandlers[s.nextIOHandler]
				s.nextIOHandler = (s.nextIOHandler + 1) % len(s.ioHandlers)

				if err := handler.AddConn(conn); err != nil {
					log.Printf("Failed to add connection to I/O handler %d: %v", handler.id, err)
					// If adding fails, close the connection to avoid resource leak
					conn.Close()
				}
			}
		}()
	}

	return nil
}
