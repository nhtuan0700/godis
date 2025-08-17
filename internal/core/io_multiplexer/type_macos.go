//go:build darwin

package io_multiplexer

import "syscall"

// get syscall Event from Operation of generic Event
func (e Event) toNative(flags uint16) syscall.Kevent_t {
	var filter int16 = syscall.EVFILT_READ // read event
	if e.Op == OpWrite {
		filter = syscall.EVFILT_WRITE
	}

	return syscall.Kevent_t{
		Ident:  uint64(e.Fd),
		Filter: filter,
		Flags:  flags,
	}
}

func createEvent(kEvent syscall.Kevent_t) Event {
	var op Operation = OpRead
	if kEvent.Filter == syscall.EVFILT_WRITE {
		op = OpWrite
	}
	return Event{
		Fd: int(kEvent.Ident),
		Op: op,
	}
}
