package core

import "errors"

const CRLF = "\r\n"

var RespNil = "$-1\r\n"

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

func readInt64(data []byte) (int64, int, error) {
	return 0, 0, errors.New("to be implemented")
}

func readBulkString(data []byte) (string, int, error) {
	return "", 0, errors.New("to be implemented")
}

func readArray(data []byte) (any, int, error) {
	return nil, 0, errors.New("to be implemented")
}

func readError(data []byte) (string, int, error) {
	return "", 0, errors.New("to be implemented")
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

// convert from resp's data to raw data
func Decode(data []byte) (any, error) {
	res, _, err := DecodeOne(data)
	return res, err
}

// convert from raw data to resp's data
func Encode(value any) []byte {
	return nil
}
