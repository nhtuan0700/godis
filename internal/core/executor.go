package core

import (
	"errors"
	"fmt"
	"strconv"
	"syscall"
	"time"

	"github.com/nhtuan0700/godis/internal/constant"
)

// PING [message]
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

// SET key value [EX seconds]
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

// GET key
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

// TTL key
func cmdTTL(args []string) []byte {
	if len(args) > 1 {
		return Encode(errors.New("ERR wrong number of arguments for 'ttl' command"), false)
	}

	key := args[0]
	obj := dictStore.Get(key)
	if obj == nil {
		return constant.RespKeyNotExist
	}

	exp, ok := dictStore.GetExpiry(key)
	if !ok {
		return constant.TTLKeyExistNoExpire
	}
	remains := int64(exp - uint64(time.Now().UnixMilli()))
	if remains < 0 {
		return constant.RespKeyNotExist
	}

	return Encode(remains/1000, false)
}

// PTTL key
func cmdPTTL(args []string) []byte {
	if len(args) > 1 {
		return Encode(errors.New("ERR wrong number of arguments for 'pttl' command"), false)
	}

	key := args[0]
	obj := dictStore.Get(key)
	if obj == nil {
		return constant.RespKeyNotExist
	}

	exp, ok := dictStore.GetExpiry(key)
	if !ok {
		return constant.TTLKeyExistNoExpire
	}
	remains := int64(exp - uint64(time.Now().UnixMilli()))
	if remains < 0 {
		return constant.RespKeyNotExist
	}

	return Encode(remains, false)
}

// DEL key [key ...]
func cmdDel(args []string) []byte {
	if len(args) == 0 {
		return Encode(errors.New("ERR wrong number of arguments for 'del' command"), false)
	}

	delCount := 0
	for _, key := range args {
		obj := dictStore.Get(key)
		if obj == nil {
			continue
		}
		delCount++
	}

	return Encode(delCount, false)
}

// EXISTS key [key ...]
func cmdExists(args []string) []byte {
	if len(args) == 0 {
		return Encode(errors.New("ERR wrong number of arguments for 'exists' command"), false)
	}

	existingCount := 0
	for _, key := range args {
		obj := dictStore.Get(key)
		if obj == nil {
			continue
		}
		existingCount++
	}

	return Encode(existingCount, false)
}

// EXPIRE key seconds
func cmdExpire(args []string) []byte {
	if len(args) != 2 {
		return Encode(errors.New("ERR wrong number of arguments for 'expire' command"), false)
	}

	key, value := args[0], args[1]
	obj := dictStore.Get(key)
	if obj == nil {
		return Encode(0, false)
	}

	expiredSec, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return Encode(errors.New("ERR value is not an integer or out of range"), false)
	}
	if expiredSec <= 0 {
		dictStore.Del(key)
	} else {
		dictStore.SetExpiry(key, uint64(expiredSec*1000))
	}

	return Encode(1, false)
}

// ExecuteAndResponse given a command, executes it and response
func ExecuteAndResponse(cmd *Command, connFd int) error {
	var res []byte

	switch cmd.Cmd {
	case constant.CMD_PING:
		res = cmdPING(cmd.Args)
	case constant.CMD_SET:
		res = cmdSet(cmd.Args)
	case constant.CMD_GET:
		res = cmdGet(cmd.Args)
	case constant.CMD_TTL:
		res = cmdTTL(cmd.Args)
	case constant.CMD_PTTL:
		res = cmdPTTL(cmd.Args)
	case constant.CMD_DEL:
		res = cmdDel(cmd.Args)
	case constant.CMD_EXIST:
		res = cmdExists(cmd.Args)
	case constant.CMD_EXPIRE:
		res = cmdExpire(cmd.Args)
	case constant.CMD_SADD:
		res = cmdSADD(cmd.Args)
	case constant.CMD_SREM:
		res = cmdSREM(cmd.Args)
	case constant.CMD_SISMEMBER:
		res = cmdSISMEMBER(cmd.Args)
	case constant.CMD_SMEMBERS:
		res = cmdSMEMEBERS(cmd.Args)
	case constant.CMD_ZADD:
		res = cmdZADD(cmd.Args)
	case constant.CMD_ZSCORE:
		res = cmdZSCORE(cmd.Args)
	case constant.CMD_ZRANK:
		res = cmdZRANK(cmd.Args)
	default:
		res = []byte(fmt.Sprintf("-ERR unknown command %s, with args beginning with:\r\n", cmd.Cmd))
	}

	_, err := syscall.Write(connFd, res)
	return err
}
