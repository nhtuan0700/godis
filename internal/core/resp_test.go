package core_test

import (
	"fmt"
	"testing"

	"github.com/nhtuan0700/godis/internal/core"
	"github.com/stretchr/testify/assert"
)

const ErrMsgIncorrectRESP = "incorrect RESP standard format"

func TestMain(m *testing.M) {
	m.Run()
}

func TestSimpleStringDecode(t *testing.T) {
	testCases := []struct {
		name     string
		respData string
		expected string
		errorMsg string
	}{
		{
			name:     "valid data",
			respData: "+OK\r\n",
			expected: "OK",
		},
		{
			name:     "incorrect RESP format",
			respData: "+OK",
			expected: "",
			errorMsg: ErrMsgIncorrectRESP,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			value, err := core.Decode([]byte(tt.respData))
			if tt.errorMsg != "" {
				assert.Error(t, err, tt.errorMsg)
			}
			assert.Equal(t, tt.expected, value)
		})
	}
}

func TestInt64Decode(t *testing.T) {
	testCases := []struct {
		name     string
		respData string
		expected int64
		errorMsg string
	}{
		{
			name:     "valid data",
			respData: ":0\r\n",
			expected: 0,
		},
		{
			name:     "valid data",
			respData: ":100\r\n",
			expected: 100,
		},
		{
			name:     "valid data",
			respData: ":-100\r\n",
			expected: -100,
		},
		{
			name:     "valid data",
			respData: ":+100\r\n",
			expected: 100,
		},
		{
			name:     "incorrect RESP format",
			respData: ":-100",
			errorMsg: ErrMsgIncorrectRESP,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			value, err := core.Decode([]byte(tt.respData))
			if tt.errorMsg != "" {
				assert.Equal(t, err.Error(), tt.errorMsg)
			}
			assert.Equal(t, tt.expected, value)
		})
	}
}

func TestBulkStringDecode(t *testing.T) {
	testCases := []struct {
		name     string
		respData string
		expected string
		errorMsg string
	}{
		{
			name:     "valid data",
			respData: "$5\r\nhello\r\n",
			expected: "hello",
		},
		{
			name:     "valid data",
			respData: "$0\r\n\r\n",
			expected: "",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			value, err := core.Decode([]byte(tt.respData))
			if tt.errorMsg != "" {
				assert.Equal(t, err.Error(), tt.errorMsg)
			}
			assert.Equal(t, tt.expected, value)
		})
	}
}

func TestArrayDecode(t *testing.T) {
	testCases := []struct {
		name     string
		respData string
		expected []any
		errorMsg string
	}{
		{
			name:     "valid data",
			respData: "*0\r\n",
			expected: []any{},
		},
		{
			name:     "valid data",
			respData: "*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n",
			expected: []any{"hello", "world"},
		},
		{
			name:     "valid data",
			respData: "*3\r\n:1\r\n:2\r\n:3\r\n",
			expected: []any{int64(1), int64(2), int64(3)},
		},
		{
			name:     "valid data",
			respData: "*5\r\n:1\r\n:2\r\n:3\r\n:4\r\n$5\r\nhello\r\n",
			expected: []any{int64(1), int64(2), int64(3), int64(4), "hello"},
		},
		{
			name:     "valid data",
			respData: "*2\r\n*3\r\n:1\r\n:2\r\n:3\r\n*2\r\n+Hello\r\n-World\r\n",
			expected: []any{[]int64{int64(1), int64(2), int64(3)}, []any{"Hello", "World"}},
		},
		{
			name:     "valid data",
			respData: "*3\r\n:1\r\n:2\r\n:3\r\n",
			expected: []any{int64(1), int64(2), int64(3)},
		},
	}

	for _, tt := range testCases {
		value, _ := core.Decode([]byte(tt.respData))
		array := value.([]any)
		assert.Equal(t, len(array), len(tt.expected))
		for i := range array {
			assert.Equal(t, fmt.Sprintf("%v", tt.expected[i]), fmt.Sprintf("%v", array[i]))
		}
	}
}

func TestEncodeString2DArray(t *testing.T) {
	var decode = [][]string{{"hello", "world"}, {"1", "2", "3"}, {"xyz"}}
	encode := core.Encode(decode, false)
	assert.EqualValues(t, "*3\r\n*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n*3\r\n$1\r\n1\r\n$1\r\n2\r\n$1\r\n3\r\n*1\r\n$3\r\nxyz\r\n", string(encode))
	decodeAgain, _ := core.Decode(encode)
	for i := 0; i < 3; i++ {
		for j := 0; j < len(decode[i]); j++ {
			assert.EqualValues(t, decode[i][j], decodeAgain.([]interface{})[i].([]interface{})[j])
		}
	}
}
