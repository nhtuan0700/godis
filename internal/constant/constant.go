package constant

const CRLF = "\r\n"

var (
	RespNil                 = []byte("$-1\r\n")
	RespOk                  = []byte("+OK\r\n")
	RespKeyNotExist         = []byte(":-2\r\n")
	TTLKeyExistNoExpire     = []byte(":-1\r\n")
	ActiveExpireSampleSized = 20
	ActiveExpireThreshold   = 0.1
)
