package core

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/nhtuan0700/godis/internal/constant"
	"github.com/nhtuan0700/godis/internal/core/data_structure"
)

// ZADD key score member [score member ...]
func cmdZADD(args []string) []byte {
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
	zset, exist := zsetStore[key]
	if !exist {
		zset = data_structure.CreateZSet(constant.BPlusTreeDegree)
		zsetStore[key] = zset
	}

	added := 0
	for i := startScoreIdx; i < len(args); i += 2 {
		scoreStr, member := args[i], args[i+1]
		scoreNum, err := strconv.ParseFloat(scoreStr, 64)
		if err != nil {
			return Encode(errors.New("score must be floating point number"), false)
		}
		ret := zset.Add(scoreNum, member)
		if ret == 1 {
			added++
		}
	}

	return Encode(added, false)
}

// ZSCORE key member
func cmdZSCORE(args []string) []byte {
	if len(args) != 2 {
		return Encode(errors.New("ERR wrong number of arguments for 'zscore' command"), false)
	}

	key, member := args[0], args[1]
	zset, exist := zsetStore[key]
	if !exist {
		return constant.RespNil
	}

	score, exist := zset.GetScore(member)
	if !exist {
		return constant.RespNil
	}

	return Encode(fmt.Sprintf("%g", score), false)
}

// ZRANK key member
func cmdZRANK(args []string) []byte {
	if len(args) != 2 {
		return Encode(errors.New("ERR wrong number of arguments for 'zrank' command"), false)
	}

	key, member := args[0], args[1]
	zset, exist := zsetStore[key]
	if !exist {
		return constant.RespNil
	}

	rank := zset.GetRank(member)
	if rank == -1 {
		return constant.RespNil
	}
	return Encode(rank, false)
}
