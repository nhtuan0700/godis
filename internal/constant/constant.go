package constant

const CRLF = "\r\n"

var (
	RespNil                 = []byte("$-1\r\n")
	RespOk                  = []byte("+OK\r\n")
	RespKeyNotExist         = []byte(":-2\r\n")
	RespZero                = []byte(":0\r\n")
	RespOne                 = []byte(":1\r\n")
	TTLKeyExistNoExpire     = []byte(":-1\r\n")
	ErrorWrongTypeKey       = []byte("-WRONGTYPE Operation against a key holding the wrong kind of value\r\n")
	ActiveExpireSampleSized = 20
	ActiveExpireThreshold   = 0.1
	BPlusTreeDegree         = 4
)

const (
	DictType = iota
	SimpleSetType
)

const (
	BfDefaultInitCapacity uint64  = 100
	BfDefaultErrRate      float64 = 0.01
)

const (
	ServerStatusIdle = iota
	ServerStatusBusy
	ServerStatusShuttingDown
)
