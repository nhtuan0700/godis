package core

import (
	"errors"
	"fmt"
	"strconv"
	"syscall"

	"github.com/nhtuan0700/godis/internal/constant"
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

func cmdSet(args []string) []byte {
	if len(args) == 1 || len(args) == 3 || len(args) > 4 {
		return Encode(errors.New("ERR wrong number of arguments for 'set' command"), false)
	}

	var ttlMs uint64 = 0
	key, value := args[0], args[1]
	if len(args) > 2 {
		ttlSec, err := strconv.ParseInt(args[3], 10, 64)
		if err != nil {
			return Encode(errors.New("ERR value is not an integer or out of range"), false)
		}
		ttlMs = uint64(ttlSec) * 1000
	}

	dictStore.Set(key, dictStore.NewObj(key, value, ttlMs))
	return constant.RespOk
}

func cmdGet(args []string) []byte {
	if len(args) != 1 {
		return Encode(errors.New("ERR wrong number of arguments for 'get' command"), false)
	}

	key := args[0]
	obj := dictStore.Get(key)
	if obj == nil {
		return constant.RespNil
	}

	if dictStore.HasExpired(key) {
		return constant.RespNil
	}

	return Encode(obj.Value, false)
}

func ExcuteAndResponse(cmd *Command, connFd int) error {
	var res []byte

	switch cmd.Cmd {
	case constant.CMD_PING:
		res = cmdPING(cmd.Args)
	case constant.CMD_SET:
		res = cmdSet(cmd.Args)
	case constant.CMD_GET:
		res = cmdGet(cmd.Args)
	default:
		res = []byte(fmt.Sprintf("-ERR unknown command %s, with args beginning with:\r\n", cmd.Cmd))
	}

	_, err := syscall.Write(connFd, res)
	return err
}
