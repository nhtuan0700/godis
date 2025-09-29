package core

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/nhtuan0700/godis/internal/constant"
)

// +OK\r\n => OK, 5, nil
// return params:
// 1: original string
// 2: next pos for next string
// 3: error
func readSimpleString(data []byte) (string, int, error) {
	pos := 1
	for pos < len(data) && data[pos] != '\r' {
		pos++
	}
	if pos == len(data) {
		return "", 0, errors.New("incorrect RESP standard format")
	}
	return string(data[1:pos]), pos + 2, nil
}

// :-100\r\n => -100, 7, nil
func readInt64(data []byte) (int64, int, error) {
	pos := 1
	var signed int64 = 1
	if data[pos] == '-' {
		signed = -1
		pos++
	}
	if data[pos] == '+' {
		pos++
	}

	var value int64
	for pos < len(data) && data[pos] != '\r' {
		value = value*10 + int64(data[pos]-'0')
		pos++
	}

	if pos == len(data) {
		return 0, 0, errors.New("incorrect RESP standard format")
	}
	return signed * value, pos + 2, nil
}

// $5\r\nhello\r\n => 5, 4
func readLen(data []byte) (int, int) {
	res, pos, _ := readInt64(data)
	return int(res), pos
}

// $5\r\nhello\r\n => hello, 11
func readBulkString(data []byte) (string, int, error) {
	length, pos := readLen(data)
	return string(data[pos:(pos + length)]), pos + length + 2, nil
}

// *2\r\n$5\r\nhello\r\n$5\r\nworld\r\n => {"hello", "world"}
func readArray(data []byte) (any, int, error) {
	length, pos := readLen(data)

	res := make([]any, length)
	for i := 0; i < length; i++ {
		elm, delta, err := DecodeOne(data[pos:])
		if err != nil {
			return nil, 0, err
		}
		res[i] = elm
		pos += delta
	}

	return res, pos, nil
}

func readError(data []byte) (string, int, error) {
	return readSimpleString(data)
}

func DecodeOne(data []byte) (any, int, error) {
	if len(data) == 0 {
		return nil, 0, errors.New("no data")
	}

	switch data[0] {
	case '+':
		return readSimpleString(data)
	case '$':
		return readBulkString(data)
	case ':':
		return readInt64(data)
	case '*':
		return readArray(data)
	case '-':
		return readError(data)
	}

	return nil, 0, nil
}

// RESP data => raw data
func Decode(data []byte) (any, error) {
	res, _, err := DecodeOne(data)
	return res, err
}

func encodeString(s string) []byte {
	return []byte(fmt.Sprintf("$%d%s%s%s", len(s), constant.CRLF, s, constant.CRLF))
}

func encodeStringArray(sa []string) []byte {
	var b []byte
	buf := bytes.NewBuffer(b)
	for _, s := range sa {
		buf.Write(encodeString(s))
	}

	return []byte(fmt.Sprintf("*%d%s%s", len(sa), constant.CRLF, buf.Bytes()))
}

// raw data => RESP data
func Encode(value any, isSimpleString bool) []byte {
	switch v := value.(type) {
	case string:
		if isSimpleString {
			return []byte(fmt.Sprintf("+%s%s", v, constant.CRLF))
		}
		return []byte(fmt.Sprintf("$%d%s%s%s", len(v), constant.CRLF, v, constant.CRLF))
	case int, int8, int16, int32, int64, uint, uint8, uint32, uint64:
		return []byte(fmt.Sprintf(":%d%s", v, constant.CRLF))
	case error:
		return []byte(fmt.Sprintf("-%s%s", v.Error(), constant.CRLF))
	case []string:
		return encodeStringArray(v)
	case [][]string:
		var b []byte
		buf := bytes.NewBuffer(b)
		for _, x := range v {
			buf.Write(encodeStringArray(x))
		}
		return []byte(fmt.Sprintf("*%d%s%s", len(v), constant.CRLF, buf.Bytes()))
	case []any:
		var b []byte
		buf := bytes.NewBuffer(b)
		for _, x := range v {
			buf.Write(Encode(x, false))
		}
		return []byte(fmt.Sprintf("*%d%s%s", len(v), constant.CRLF, buf.Bytes()))

	default:
		return constant.RespNil
	}
}

func ParseCommand(data []byte) (*Command, error) {
	value, err := Decode(data)
	if err != nil {
		return nil, err
	}

	array := value.([]any)
	tokens := make([]string, len(array))
	for i := range tokens {
		tokens[i] = array[i].(string)
	}

	return &Command{Cmd: strings.ToUpper(tokens[0]), Args: tokens[1:]}, nil
}
