package data_structure

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateCMSByDIM(t *testing.T) {
	width, depth := 1, 1

	cms := CreateCMS(uint32(width), uint32(depth))
	assert.EqualValues(t, cms.width, width)
	assert.EqualValues(t, cms.depth, depth)
}

func TestCalcCMSDim(t *testing.T) {
	testCases := []struct {
		errRate       float64
		probability   float64
		expectedWidth uint32
		expectedDepth uint32
	}{
		{
			errRate:       0.01,
			probability:   0.01,
			expectedWidth: 200,
			expectedDepth: 7,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TestCalcCMSDim#%d", i), func(t *testing.T) {
			width, depth := CalcCMSDim(tc.errRate, tc.probability)
			assert.EqualValues(t, tc.expectedWidth, width)
			assert.EqualValues(t, tc.expectedDepth, depth)
		})
	}
}

func TestIncrBy_DIM(t *testing.T) {
	type arg struct {
		key   string
		count uint64
	}

	testCases := []struct {
		name     string
		width    uint32
		depth    uint32
		list     []arg
		expected []arg
	}{
		{
			name:  "Collision",
			width: 1,
			depth: 1,
			list:  []arg{{"a", 1}, {"b", 1}, {"a", 1}, {"b", 1}},
			expected: []arg{
				{"a", 4},
				{"b", 4},
			},
		},
		{
			name:  "Not Collision",
			width: 200,
			depth: 8,
			list:  []arg{{"a", 1}, {"b", 1}, {"a", 1}, {"b", 1}},
			expected: []arg{
				{"a", 2},
				{"b", 2},
				{"notknown", 0},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cms := CreateCMS(tc.width, tc.depth)

			for i := 0; i < len(tc.list); i++ {
				key, count := tc.list[i].key, tc.list[i].count
				cms.IncrBy(key, count)
			}

			for i := 0; i < len(tc.expected); i++ {
				assert.EqualValues(t, tc.expected[i].count, cms.Count(tc.expected[i].key))
			}
		})
	}
}

func TestIncrBy_PROB(t *testing.T) {
	type arg struct {
		key   string
		count uint64
	}

	testCases := []struct {
		name        string
		errRate     float64
		probability float64
		list        []arg
		expected    []arg
	}{
		{
			name:        "Not Collision",
			errRate:     0.01,
			probability: 0.01,
			list:        []arg{{"a", 1}, {"b", 1}, {"a", 1}, {"b", 1}},
			expected: []arg{
				{"a", 2},
				{"b", 2},
				{"notknown", 0},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			width, depth := CalcCMSDim(tc.errRate, tc.probability)
			cms := CreateCMS(width, depth)

			for i := 0; i < len(tc.list); i++ {
				key, count := tc.list[i].key, tc.list[i].count
				cms.IncrBy(key, count)
			}

			for i := 0; i < len(tc.expected); i++ {
				assert.EqualValues(t, tc.expected[i].count, cms.Count(tc.expected[i].key))
			}
		})
	}
}

func TestCount(t *testing.T) {
	cms := CreateCMS(200, 7)
	type arg struct {
		key   string
		count uint64
	}

	list := []arg{
		{"a", 3},
		{"b", 4},
		{"c", 5},
		{"a", 3},
		{"b", 4},
		{"c", 5},
	}

	for i := 0; i < len(list); i++ {
		cms.IncrBy(list[i].key, list[i].count)
	}

	assert.EqualValues(t, 6, cms.Count("a"))
	assert.EqualValues(t, 8, cms.Count("b"))
	assert.EqualValues(t, 10, cms.Count("c"))
}
