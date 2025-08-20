package core

import (
	"errors"
	"fmt"
	"syscall"
)

func cmdPING(args []string) []byte {
	switch {
	case len(args) == 0:
		return Encode("PONG", true)
	case len(args) == 1:
		return Encode(args[0], false)
	default:
		return Encode(errors.New("ERR wrong number of arguments for 'ping' command"), false)
	}
}

func ExcuteAndResponse(cmd *Command, connFd int) error {
	var res []byte

	switch cmd.Cmd {
	case CMD_PING:
		res = cmdPING(cmd.Args)
	default:
		res = []byte(fmt.Sprintf("-ERR unknown command %s, with args beginning with:\r\n", cmd.Cmd))
	}

	_, err := syscall.Write(connFd, res)
	return err
}
