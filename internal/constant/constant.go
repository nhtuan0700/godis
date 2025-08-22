package constant

const CRLF = "\r\n"

var RespNil = []byte("$-1\r\n")
var RespOk = []byte("+OK\r\n")
var RespKeyNotExist = []byte(":-2\r\n")
var TTLKeyExistNoExpire = []byte(":-1\r\n")
