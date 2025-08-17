//go:build linux

package io_multiplexer

import "syscall"

func (e Event) toNative() syscall.EpollEvent {
	var event uint32 = syscall.EPOLLIN
	if e.Op == OpWrite {
		event = syscall.EPOLLOUT
	}

	return syscall.EpollEvent{
		Fd: int32(e.Fd),
		Events: event,
	}
}

func createEvent(epollEvent syscall.EpollEvent) Event {
	var op Operation = OpRead
	if epollEvent.Events == syscall.EPOLLOUT {
		op = OpWrite
	}

	return Event{
		Fd: int(epollEvent.Fd),
		Op: op,
	}
}
