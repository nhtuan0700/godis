package core

import (
	"errors"
	"strconv"

	"github.com/nhtuan0700/godis/internal/constant"
	"github.com/nhtuan0700/godis/internal/core/data_structure"
)

// CMS.INITBYDIM key width depth
func cmdCMSINITBYDIM(redisDB *RedisDB, args []string) []byte {
	if len(args) != 3 {
		return Encode(errors.New("ERR wrong number of arguments for 'cms.initbydim' command"), false)
	}

	key := args[0]
	obj, exist := redisDB.dict[key]
	if exist {
		_, ok := obj.value.(*data_structure.CMS)
		if !ok {
			return constant.ErrorWrongTypeKey
		}
		return Encode(errors.New("CMS: key already exists"), false)
	}

	width, err := strconv.ParseUint(args[1], 10, 32)
	if err != nil {
		return Encode(errors.New("CMS: invalid width"), false)
	}
	depth, err := strconv.ParseUint(args[2], 10, 32)
	if err != nil {
		return Encode(errors.New("CMS: invalid depth"), false)
	}

	cms := data_structure.CreateCMS(uint32(width), uint32(depth))
	redisDB.dict[key] = NewRedisObj(cms)

	return constant.RespOk
}

// CMS.INITBYPROB key error probability
func cmdCMSINITBYPROB(redisDB *RedisDB, args []string) []byte {
	if len(args) != 3 {
		return Encode(errors.New("ERR wrong number of arguments for 'cms.initbyprob' command"), false)
	}
	key := args[0]
	obj, exist := redisDB.dict[key]
	if exist {
		_, ok := obj.value.(*data_structure.CMS)
		if !ok {
			return constant.ErrorWrongTypeKey
		}
		return Encode(errors.New("CMS: key already exists"), false)
	}

	errRate, err := strconv.ParseFloat(args[1], 64)
	if err != nil || errRate <= 0 || errRate >= 1 {
		return Encode(errors.New("CMS: invalid overestimation value"), false)
	}
	probability, err := strconv.ParseFloat(args[2], 64)
	if err != nil || probability <= 0 || probability >= 1 {
		return Encode(errors.New("CMS: invalid prob value"), false)
	}

	width, depth := data_structure.CalcCMSDim(errRate, probability)
	cms := data_structure.CreateCMS(width, depth)
	redisDB.dict[key] = NewRedisObj(cms)

	return constant.RespOk
}

// CMS.INCRBY key item value
func cmdCMSINCRBY(redisDB *RedisDB, args []string) []byte {
	if len(args) < 3 || len(args)%2 == 0 {
		return Encode(errors.New("ERR wrong number of arguments for 'cms.incrby' command"), false)
	}

	key := args[0]
	obj, exist := redisDB.dict[key]
	if !exist {
		return Encode(errors.New("CMS: key does not exist"), false)
	}
	cms, ok := obj.value.(*data_structure.CMS)
	if !ok {
		return constant.ErrorWrongTypeKey
	}

	res := make([]any, 0)
	for i := 1; i < len(args); i += 2 {
		item := args[i]
		count, err := strconv.ParseUint(args[i+1], 10, 64)
		if err != nil {
			return Encode(errors.New("CMS: Cannot parse number"), false)
		}

		res = append(res, cms.IncrBy(item, uint64(count)))
	}

	return Encode(res, false)
}

// CMS.QUERY key item [item ...]
func cmdCMSQUERY(redisDB *RedisDB, args []string) []byte {
	if len(args) < 2 {
		return Encode(errors.New("ERR wrong number of arguments for 'cms.query' command"), false)
	}

	key := args[0]
	obj, exist := redisDB.dict[key]
	if !exist {
		return Encode(errors.New("CMS: key does not exist"), false)
	}
	cms, ok := obj.value.(*data_structure.CMS)
	if !ok {
		return constant.ErrorWrongTypeKey
	}

	res := make([]any, 0)
	for i := 1; i < len(args); i++ {
		res = append(res, cms.Count(args[i]))
	}

	return Encode(res, false)
}
