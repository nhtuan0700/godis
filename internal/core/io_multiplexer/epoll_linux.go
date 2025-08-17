//go:build linux

package io_multiplexer

import (
	"syscall"

	"github.com/nhtuan0700/godis/internal/config"
)

type Epoll struct {
	fd            int
	epollEvents   []syscall.EpollEvent // temporary buffer
	genericEvents []Event
}

func CreateIOMultiplexer() (*Epoll, error) {
	epollFD, err := syscall.EpollCreate1(0)
	if err != nil {
		return nil, err
	}

	return &Epoll{
		fd:            epollFD,
		epollEvents:   make([]syscall.EpollEvent, config.MaxConnections),
		genericEvents: make([]Event, config.MaxConnections),
	}, nil
}

// Subscribe file descriptor's event to the monitoring list
func (ep *Epoll) Monitor(event Event) error {
	epollEvent := event.toNative()
	// add event.fd to the monitoring list of ep.fd
	return syscall.EpollCtl(ep.fd, syscall.EPOLL_CTL_ADD, event.Fd, &epollEvent)
}

// Wait for events in the monitoring list
func (ep *Epoll) Wait() ([]Event, error) {
	n, err := syscall.EpollWait(ep.fd, ep.epollEvents, -1)
	if err != nil {
		return nil, err
	}

	for i := 0; i < n; i++ {
		ep.genericEvents[i] = createEvent(ep.epollEvents[i])
	}

	return ep.genericEvents[:n], nil
}

func (ep *Epoll) Close() error {
	return syscall.Close(ep.fd)
}
