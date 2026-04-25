package core

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/nhtuan0700/godis/internal/constant"
)

type Command struct {
	Cmd  string
	Args []string
}

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
func cmdSet(redisDB *RedisDB, args []string) []byte {
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

	redisDB.Set(key, NewRedisObj(value), ttlMs)
	return constant.RespOk
}

// GET key
func cmdGet(redisDB *RedisDB, args []string) []byte {
	if len(args) != 1 {
		return Encode(errors.New("ERR wrong number of arguments for 'get' command"), false)
	}

	key := args[0]
	obj := redisDB.Get(key)
	if obj == nil {
		return constant.RespNil
	}

	if redisDB.HasExpired(key) {
		return constant.RespNil
	}

	return Encode(obj.value, false)
}

// TTL key
func cmdTTL(redisDB *RedisDB, args []string) []byte {
	if len(args) > 1 {
		return Encode(errors.New("ERR wrong number of arguments for 'ttl' command"), false)
	}

	key := args[0]
	obj := redisDB.Get(key)
	if obj == nil {
		return constant.RespKeyNotExist
	}

	exp, ok := redisDB.GetExpiry(key)
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
func cmdPTTL(redisDB *RedisDB, args []string) []byte {
	if len(args) > 1 {
		return Encode(errors.New("ERR wrong number of arguments for 'pttl' command"), false)
	}

	key := args[0]
	obj := redisDB.Get(key)
	if obj == nil {
		return constant.RespKeyNotExist
	}

	exp, ok := redisDB.GetExpiry(key)
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
func cmdDel(redisDB *RedisDB, args []string) []byte {
	if len(args) == 0 {
		return Encode(errors.New("ERR wrong number of arguments for 'del' command"), false)
	}

	delCount := 0
	for _, key := range args {
		ok := redisDB.Delete(key)
		if ok {
			delCount++
		}
	}

	return Encode(delCount, false)
}

// EXISTS key [key ...]
func cmdExists(redisDB *RedisDB, args []string) []byte {
	if len(args) == 0 {
		return Encode(errors.New("ERR wrong number of arguments for 'exists' command"), false)
	}

	existingCount := 0
	for _, key := range args {
		obj := redisDB.Get(key)
		if obj == nil {
			continue
		}
		existingCount++
	}

	return Encode(existingCount, false)
}

// EXPIRE key seconds
func cmdExpire(redisDB *RedisDB, args []string) []byte {
	if len(args) != 2 {
		return Encode(errors.New("ERR wrong number of arguments for 'expire' command"), false)
	}

	key, value := args[0], args[1]
	obj := redisDB.Get(key)
	if obj == nil {
		return Encode(0, false)
	}

	expiredSec, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return Encode(errors.New("ERR value is not an integer or out of range"), false)
	}
	if expiredSec <= 0 {
		redisDB.Delete(key)
	} else {
		redisDB.SetExpiry(key, uint64(expiredSec*1000))
	}

	return Encode(1, false)
}

// INFO [section [section...]]
func cmdINFO(redisDB *RedisDB, args []string) []byte {
	var info []byte
	buf := bytes.NewBuffer(info)
	buf.WriteString("# Keyspace\r\n")
	buf.WriteString(fmt.Sprintf("db:key=%d,epxires=%d,avg_ttl=0\r\n", len(redisDB.dict), len(redisDB.expireDict)))
	return Encode(buf.String(), false)
}
