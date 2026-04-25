package core

import (
	"errors"

	"github.com/nhtuan0700/godis/internal/constant"
	"github.com/nhtuan0700/godis/internal/core/data_structure"
)

// SADD key member [member ...]
func cmdSADD(redisDB *RedisDB, args []string) []byte {
	if len(args) <= 1 {
		return Encode(errors.New("ERR wrong number of arguments for 'sadd' command"), false)
	}

	key := args[0]
	var simpleSet *data_structure.SimpleSet
	obj, exist := redisDB.dict[key]
	if !exist {
		simpleSet = data_structure.NewSimpleSet()
		redisDB.dict[key] = NewRedisObj(simpleSet)
		return Encode(simpleSet.Add(args[1:]...), false)
	}
	simpleSet, ok := obj.value.(*data_structure.SimpleSet)
	if !ok {
		return constant.ErrorWrongTypeKey
	}
	return Encode(simpleSet.Add(args[1:]...), false)
}

// SREM key member [member ...]
func cmdSREM(redisDB *RedisDB, args []string) []byte {
	if len(args) <= 1 {
		return Encode(errors.New("ERR wrong number of arguments for 'srem' command"), false)
	}

	key := args[0]
	obj, exist := redisDB.dict[key]
	if !exist {
		return Encode(0, false)
	}
	simpleSet, ok := obj.value.(*data_structure.SimpleSet)
	if !ok {
		return constant.ErrorWrongTypeKey
	}

	return Encode(simpleSet.Remove(args[1:]...), false)
}

// SISMEMBER key member
func cmdSISMEMBER(redisDB *RedisDB, args []string) []byte {
	if len(args) != 2 {
		return Encode(errors.New("ERR wrong number of arguments for 'sismember' command"), false)
	}
	key, member := args[0], args[1]
	obj, exist := redisDB.dict[key]
	if !exist {
		return Encode(0, false)
	}
	simpleSet, ok := obj.value.(*data_structure.SimpleSet)
	if !ok {
		return constant.ErrorWrongTypeKey
	}

	return Encode(simpleSet.IsMember(member), false)
}

// SMEMBERS key
func cmdSMEMEBERS(redisDB *RedisDB, args []string) []byte {
	if len(args) != 1 {
		return Encode(errors.New("ERR wrong number of arguments for 'smembers' command"), false)
	}

	key := args[0]
	obj, exist := redisDB.dict[key]
	if !exist {
		return Encode(make([]string, 0), false)
	}
	simpleSet, ok := obj.value.(*data_structure.SimpleSet)
	if !ok {
		return constant.ErrorWrongTypeKey
	}

	return Encode(simpleSet.Members(), false)
}
