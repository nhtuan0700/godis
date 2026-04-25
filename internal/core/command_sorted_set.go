package core

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/nhtuan0700/godis/internal/constant"
	"github.com/nhtuan0700/godis/internal/core/data_structure"
)

// ZADD key score member [score member ...]
func cmdZADD(redisDB *RedisDB, args []string) []byte {
	if len(args) < 3 {
		return Encode(errors.New("ERR wrong number of arguments for 'zadd' command"), false)
	}

	// 0: key of sorted set
	startScoreIdx := 1
	numElems := len(args) - startScoreIdx
	if numElems%2 == 1 {
		return Encode(errors.New("ERR syntax error"), false)
	}

	key := args[0]
	var zset *data_structure.ZSet
	obj, exist := redisDB.dict[key]
	if !exist {
		zset = data_structure.NewZSet()
		redisDB.dict[key] = NewRedisObj(zset)
	} else {
		var ok bool
		zset, ok = obj.value.(*data_structure.ZSet)
		if !ok {
			return constant.ErrorWrongTypeKey
		}
	}

	added := 0
	for i := startScoreIdx; i < len(args); i += 2 {
		scoreStr, member := args[i], args[i+1]
		scoreNum, err := strconv.ParseFloat(scoreStr, 64)
		if err != nil {
			return Encode(errors.New("score must be floating point number"), false)
		}
		ret := zset.Add(scoreNum, member)
		if ret {
			added++
		}
	}

	return Encode(added, false)
}

// ZSCORE key member
func cmdZSCORE(redisDB *RedisDB, args []string) []byte {
	if len(args) != 2 {
		return Encode(errors.New("ERR wrong number of arguments for 'zscore' command"), false)
	}

	key, member := args[0], args[1]
	obj, exist := redisDB.dict[key]
	if !exist {
		return constant.RespNil
	}

	zset, ok := obj.value.(*data_structure.ZSet)
	if !ok {
		return constant.ErrorWrongTypeKey
	}

	score, exist := zset.GetScore(member)
	if !exist {
		return constant.RespNil
	}

	return Encode(fmt.Sprintf("%g", score), false)
}

// ZRANK key member
func cmdZRANK(redisDB *RedisDB, args []string) []byte {
	if len(args) != 2 {
		return Encode(errors.New("ERR wrong number of arguments for 'zrank' command"), false)
	}

	key, member := args[0], args[1]
	obj, exist := redisDB.dict[key]
	if !exist {
		return constant.RespNil
	}

	zset, ok := obj.value.(*data_structure.ZSet)
	if !ok {
		return constant.ErrorWrongTypeKey
	}

	rank, exist := zset.GetRank(member, false)
	if !exist {
		return constant.RespNil
	}

	return Encode(rank, false)
}

// ZREM key member [member ...]
func cmdZREM(redisDB *RedisDB, args []string) []byte {
	if len(args) < 2 {
		return Encode(errors.New("ERR wrong number of arguments for 'zrem' command"), false)
	}

	key := args[0]
	obj, exist := redisDB.dict[key]
	if !exist {
		return constant.RespNil
	}

	zset, ok := obj.value.(*data_structure.ZSet)
	if !ok {
		return constant.ErrorWrongTypeKey
	}

	removeCount := 0
	for i := 1; i < len(args); i++ {
		ok := zset.Remove(args[i])
		if ok {
			removeCount++
		}
	}

	return Encode(removeCount, false)
}
