package core

import (
	"errors"
	"strconv"

	"github.com/nhtuan0700/godis/internal/constant"
	"github.com/nhtuan0700/godis/internal/core/data_structure"
)

// BF.RESERVE key error_rate entries
func cmdBFRESERVE(args []string) []byte {
	if len(args) != 3 {
		return Encode(errors.New("ERR wrong number of arguments for 'bf.reserve' command"), false)
	}

	key := args[0]
	_, exist := bloomStore[key]
	if exist {
		return Encode(errors.New("ERR item exists"), false)
	}

	errorRate, err := strconv.ParseFloat(args[1], 64)
	if err != nil {
		return Encode(errors.New("ERR bad error rate"), false)

	}
	if errorRate <= 0 || errorRate >= 1 {
		return Encode(errors.New("ERR error rate must be in the range (0.000000, 1.000000)"), false)
	}

	capacity, err := strconv.ParseUint(args[2], 10, 64)
	if err != nil {
		return Encode(errors.New("ERR bad capacity"), false)
	}
	// [1, 2^30]
	if capacity < 1 || capacity > 1<<30 {
		return Encode(errors.New("ERR capacity must be in the range [1, 1073741824]"), false)
	}
	bloomStore[key] = data_structure.CreateBloomFilter(capacity, errorRate)

	return constant.RespOk
}

// BF.ADD key entry
func cmdBFADD(args []string) []byte {
	if len(args) != 2 {
		return Encode(errors.New("ERR wrong number of arguments for 'bf.add' command"), false)
	}

	key, entry := args[0], args[1]
	bloom, exist := bloomStore[key]
	if !exist {
		bloomStore[key] = data_structure.CreateBloomFilter(constant.BfDefaultInitCapacity, constant.BfDefaultErrRate)
		bloom = bloomStore[key]
	}

	if bloom.Add(entry) {
		return constant.RespOne
	}
	return constant.RespZero
}

// BF.MADD key entry [entry ...]
func cmdBFMADD(args []string) []byte {
	if len(args) < 2 {
		return Encode(errors.New("ERR wrong number of arguments for 'bf.madd' command"), false)
	}

	key := args[0]
	bloom, exist := bloomStore[key]
	if !exist {
		bloomStore[key] = data_structure.CreateBloomFilter(constant.BfDefaultInitCapacity, constant.BfDefaultErrRate)
		bloom = bloomStore[key]
	}

	res := make([]any, 0)
	for i := 1; i < len(args); i++ {
		ret := 0
		if bloom.Add(args[i]) {
			ret = 1
		}
		res = append(res, ret)
	}

	return Encode(res, false)
}

// BF.EXISTS key entry
func cmdBFEXISTS(args []string) []byte {
	if len(args) != 2 {
		return Encode(errors.New("ERR wrong number of arguments for 'bf.exists' command"), false)
	}
	key, entry := args[0], args[1]
	bloom, exist := bloomStore[key]
	if !exist {
		return Encode(constant.RespZero, false)
	}
	if bloom.Exist(entry) {
		return constant.RespOne
	}

	return constant.RespZero
}

// BF.MEXISTS key entry [entry ...]
func cmdBFMEXISTS(args []string) []byte {
	if len(args) < 2 {
		return Encode(errors.New("ERR wrong number of arguments for 'bf.exists' command"), false)
	}

	res := make([]any, 0)
	key := args[0]
	bloom, exist := bloomStore[key]
	if !exist {
		for i := 1; i < len(args); i++ {
			res = append(res, 0)
		}
		return Encode(res, false)
	}

	for i := 1; i < len(args); i++ {
		ret := 0
		if bloom.Exist(args[i]) {
			ret = 1
		}
		res = append(res, ret)
	}

	return Encode(res, false)
}
