package core

import (
	"fmt"

	"github.com/nhtuan0700/godis/internal/constant"
)

// ExecuteCommand given a command, executes it and response
func ExecuteCommand(redisDB *RedisDB, cmd *Command) []byte {
	var res []byte

	switch cmd.Cmd {
	case constant.CMD_PING:
		res = cmdPING(cmd.Args)
	case constant.CMD_SET:
		res = cmdSet(redisDB, cmd.Args)
	case constant.CMD_GET:
		res = cmdGet(redisDB, cmd.Args)
	case constant.CMD_TTL:
		res = cmdTTL(redisDB, cmd.Args)
	case constant.CMD_PTTL:
		res = cmdPTTL(redisDB, cmd.Args)
	case constant.CMD_DEL:
		res = cmdDel(redisDB, cmd.Args)
	case constant.CMD_EXIST:
		res = cmdExists(redisDB, cmd.Args)
	case constant.CMD_EXPIRE:
		res = cmdExpire(redisDB, cmd.Args)
	case constant.CMD_SADD:
		res = cmdSADD(redisDB, cmd.Args)
	case constant.CMD_SREM:
		res = cmdSREM(redisDB, cmd.Args)
	case constant.CMD_SISMEMBER:
		res = cmdSISMEMBER(redisDB, cmd.Args)
	case constant.CMD_SMEMBERS:
		res = cmdSMEMEBERS(redisDB, cmd.Args)
	case constant.CMD_ZADD:
		res = cmdZADD(redisDB, cmd.Args)
	case constant.CMD_ZSCORE:
		res = cmdZSCORE(redisDB, cmd.Args)
	case constant.CMD_ZRANK:
		res = cmdZRANK(redisDB, cmd.Args)
	case constant.CMD_ZREM:
		res = cmdZREM(redisDB, cmd.Args)
	case constant.CMD_CMS_INITBYDIM:
		res = cmdCMSINITBYDIM(redisDB, cmd.Args)
	case constant.CMD_CMS_INITBYPROB:
		res = cmdCMSINITBYPROB(redisDB, cmd.Args)
	case constant.CMD_CMS_INCRBY:
		res = cmdCMSINCRBY(redisDB, cmd.Args)
	case constant.CMD_CMS_QUERY:
		res = cmdCMSQUERY(redisDB, cmd.Args)
	case constant.CMD_BF_RESERVE:
		res = cmdBFRESERVE(redisDB, cmd.Args)
	case constant.CMD_BF_ADD:
		res = cmdBFADD(redisDB, cmd.Args)
	case constant.CMD_BF_MADD:
		res = cmdBFMADD(redisDB, cmd.Args)
	case constant.CMD_BF_EXISTS:
		res = cmdBFEXISTS(redisDB, cmd.Args)
	case constant.CMD_BF_MEXISTS:
		res = cmdBFMEXISTS(redisDB, cmd.Args)
	case constant.CMD_INFO:
		res = cmdINFO(redisDB, cmd.Args)
	default:
		res = []byte(fmt.Sprintf("-ERR unknown command %s, with args beginning with:\r\n", cmd.Cmd))
	}

	return res
}
