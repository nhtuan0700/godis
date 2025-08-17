//go:build darwin

package io_multiplexer

import (
	"syscall"

	"github.com/nhtuan0700/godis/internal/config"
)

type KQueue struct {
	fd            int
	kqEvents      []syscall.Kevent_t // temporary buffer
	genericEvents []Event
}

func CreateIOMultiplexer() (*KQueue, error) {
	kqFd, err := syscall.Kqueue()
	if err != nil {
		return nil, err
	}

	return &KQueue{
		fd:            kqFd,
		kqEvents:      make([]syscall.Kevent_t, config.MaxConnections),
		genericEvents: make([]Event, config.MaxConnections),
	}, nil
}

// Subscribe file descriptor's event to the monitoring list
func (kq *KQueue) Monitor(event Event) error {
	kqEvent := event.toNative(syscall.EV_ADD)
	// Add event.Fd to the monitoring list of kq.fd
	_, err := syscall.Kevent(kq.fd, []syscall.Kevent_t{kqEvent}, nil, nil)
	return err
}

// Wait for events in the monitoring list
func (kq *KQueue) Wait() ([]Event, error) {
	n, err := syscall.Kevent(kq.fd, nil, kq.kqEvents, nil)
	if err != nil {
		return nil, err
	}

	for i := 0; i < n; i++ {
		kq.genericEvents[i] = createEvent(kq.kqEvents[i])
	}

	return kq.genericEvents[:n], nil
}

func (kq *KQueue) Close() error {
	return syscall.Close(kq.fd)
}
