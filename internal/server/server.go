package server

import (
	"errors"
	"io"
	"log"
	"net"
	"syscall"

	"github.com/nhtuan0700/godis/internal/config"
	"github.com/nhtuan0700/godis/internal/core"
	"github.com/nhtuan0700/godis/internal/core/io_multiplexer"
)

func readCommand(fd int) (*core.Command, error) {
	var buf = make([]byte, 512)
	n, err := syscall.Read(fd, buf)
	if err != nil {
		return nil, err
	}
	if n == 0 {
		return nil, io.EOF
	}
	// log.Println("command: ", string(buf[:n]))
	return core.ParseCommand(buf)
}

func respond(cmd string, fd int) error {
	_, err := syscall.Write(fd, []byte(cmd))
	if err != nil {
		return err
	}

	return nil
}

func RunIOMultiplexingServer() error {
	// 1. Create listener FD
	listener, err := net.Listen(config.Protocol, config.Address)
	if err != nil {
		return err
	}
	defer listener.Close()
	log.Println("Starting an I/O Multiplexing TCP server on ", config.Address)

	// Get file descriptor of listener
	tcpListener, ok := listener.(*net.TCPListener)
	if !ok {
		return errors.New("listner is not a TCPListener")
	}
	listenerFile, err := tcpListener.File()
	if err != nil {
		return err
	}
	defer listenerFile.Close()

	listenerFD := int(listenerFile.Fd())

	// 2. Create an ioMultiplexer instance (epoll in Linux, kqueue in MacOS) and monitor Listner FD
	ioMultiplexer, err := io_multiplexer.CreateIOMultiplexer()
	if err != nil {
		return err
	}
	defer ioMultiplexer.Close()

	// Monitor "read" events on the Listener FD
	if err := ioMultiplexer.Monitor(io_multiplexer.Event{
		Fd: listenerFD,
		Op: io_multiplexer.OpRead,
	}); err != nil {
		return err
	}

	// 3. Monitor all the FDs in the monitoring list
	// events := make([]io_multiplexer.Event, config.MaxConnections)
	for {
		// wait for file descriptor in the monitoring list to be ready for I/O
		// it is a blocking call
		events, err := ioMultiplexer.Wait()
		if err != nil {
			continue
		}

		for i := 0; i < len(events); i++ {
			if events[i].Fd == listenerFD {
				log.Println("new client is trying to connect")
				// setup new connection
				connFd, _, err := syscall.Accept(events[i].Fd)
				if err != nil {
					log.Println("err", err)
					continue
				}
				log.Println("setup a new connection")
				// ask epoll to monitor this connection
				if err := ioMultiplexer.Monitor(io_multiplexer.Event{
					Fd: connFd,
					Op: io_multiplexer.OpRead,
				}); err != nil {
					return err
				}
			} else {
				cmd, err := readCommand(events[i].Fd)
				if err != nil {
					if err == io.EOF || err == syscall.ECONNRESET {
						log.Println("client disconnected")
						_ = syscall.Close(events[i].Fd)
						continue
					}
					log.Println("read err: ", err)
					continue
				}

				if err := core.ExcuteAndResponse(cmd, events[i].Fd); err != nil {
					log.Println("write err: ", err)
				}
			}
		}
	}
}
