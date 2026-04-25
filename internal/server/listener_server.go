package server

import (
	"context"
	"errors"
	"log"
	"net"
	"syscall"

	"github.com/nhtuan0700/godis/internal/config"
	"golang.org/x/sys/unix"
)

func (s *Server) StartSingleListener() error {
	// Start all I/O handlers event loops
	for _, handler := range s.ioHandlers {
		s.wg.Add(1)
		go func(handler *IOHandler) {
			defer s.wg.Done()
			handler.Run()
		}(handler)
	}

	// Setup listener socket

	listener, err := net.Listen(config.Protocol, config.Address)
	if err != nil {
		return err
	}
	s.addListener(listener)
	defer listener.Close()

	log.Printf("Server listening on %s", config.Address)

	for {
		if s.isDraining() {
			return nil
		}

		conn, err := listener.Accept()
		if err != nil {
			if s.isDraining() || errors.Is(err, net.ErrClosed) {
				return nil
			}
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		handler := s.nextHandler()

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
				err = syscall.SetsockoptInt(int(fd), unix.SOL_SOCKET, unix.SO_REUSEPORT, 1)
			})
			return err
		},
	}

	return lc.Listen(context.Background(), network, addr)
}

func (s *Server) StartMultiListeners() error {
	// Start all I/O handlers event loops
	for _, handler := range s.ioHandlers {
		s.wg.Add(1)
		go func(handler *IOHandler) {
			defer s.wg.Done()
			handler.Run()
		}(handler)
	}

	// Start a listener for each I/O handler
	for i := 0; i < config.ListenerNumber; i++ {
		s.wg.Add(1)
		go func(listenerID int) {
			defer s.wg.Done()

			listener, err := createReusablePortListener(config.Protocol, config.Address)
			if err != nil {
				log.Fatal(err)
			}
			s.addListener(listener)
			defer listener.Close()
			log.Printf("Listener %d started on %s", listenerID, config.Address)

			for {
				if s.isDraining() {
					return
				}

				conn, err := listener.Accept()
				if err != nil {
					if s.isDraining() || errors.Is(err, net.ErrClosed) {
						return
					}
					log.Printf("Failed to accept connection: %v", err)
					continue
				}

				handler := s.nextHandler()

				if err := handler.AddConn(conn); err != nil {
					log.Printf("Failed to add connection to I/O handler %d: %v", handler.id, err)
					// If adding fails, close the connection to avoid resource leak
					conn.Close()
				}
			}
		}(i)
	}

	return nil
}
