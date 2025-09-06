package core

import (
	"errors"

	"github.com/nhtuan0700/godis/internal/core/data_structure"
)


// SADD key member [member ...]
func cmdSADD(args []string) []byte {
	if len(args) <= 1 {
		return Encode(errors.New("ERR wrong number of arguments for 'sadd' command"), false)
	}

	key := args[0]
	simpleSet, exist := setStore[key]
	if !exist {
		simpleSet = data_structure.NewSimpleSet()
		setStore[key] = simpleSet
	}

	return Encode(simpleSet.Add(args[1:]...), false)
}

// SREM key member [member ...]
func cmdSREM(args []string) []byte {
	if len(args) <= 1 {
		return Encode(errors.New("ERR wrong number of arguments for 'srem' command"), false)
	}

	key := args[0]
	simpleSet, exist := setStore[key]
	if !exist {
		return Encode(0, false)
	}

	return Encode(simpleSet.Remove(args[1:]...), false)
}

// SISMEMBER key member
func cmdSISMEMBER(args []string) []byte {
	if len(args) != 2 {
		return Encode(errors.New("ERR wrong number of arguments for 'sismember' command"), false)
	}
	key, member := args[0], args[1]
	simpleSet, exist := setStore[key]
	if !exist {
		return Encode(0, false)
	}

	return Encode(simpleSet.IsMember(member), false)
}

// SMEMBERS key
func cmdSMEMEBERS(args []string) []byte {
	if len(args) != 1 {
		return Encode(errors.New("ERR wrong number of arguments for 'smembers' command"), false)
	}

	key := args[0]
	simpleSet, exist := setStore[key]
	if !exist {
		return Encode(make([]string, 0), false)
	}

	return Encode(simpleSet.Members(), false)
}
