package core_test

import (
	"testing"

	"github.com/nhtuan0700/godis/internal/core"
	"github.com/stretchr/testify/assert"
)


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
			errorMsg: "incorrect RESP standard format",
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
