package data_structure_test

import (
	"fmt"
	"testing"

	"github.com/nhtuan0700/godis/internal/core/data_structure"
	"github.com/stretchr/testify/assert"
)

func initSkipList() *data_structure.Skiplist {
	sl := data_structure.CreateSkiplist()

	sl.Insert(10, "k1")
	sl.Insert(40, "k4")
	sl.Insert(50, "k5")
	sl.Insert(20, "k2")
	sl.Insert(60, "k6")
	sl.Insert(80, "k8")
	sl.Insert(30, "k3")
	sl.Insert(70, "k7")

	return sl
}

func TestInitSkipList(t *testing.T) {
	testCases := []struct {
		key          string
		score        float64
		expectedRank uint32
	}{
		{
			key:          "k1",
			score:        10,
			expectedRank: 1,
		},
		{
			key:          "k2",
			score:        20,
			expectedRank: 2,
		},
		{
			key:          "k3",
			score:        30,
			expectedRank: 3,
		},
		{
			key:          "k4",
			score:        40,
			expectedRank: 4,
		},
		{
			key:          "k5",
			score:        50,
			expectedRank: 5,
		},
		{
			key:          "k6",
			score:        60,
			expectedRank: 6,
		},
		{
			key:          "k7",
			score:        70,
			expectedRank: 7,
		},
		{
			key:          "k8",
			score:        80,
			expectedRank: 8,
		},
	}

	sl := initSkipList()
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("#%d", i), func(t *testing.T) {
			rank := sl.GetRank(tc.score, tc.key)
			assert.Equal(t, rank, tc.expectedRank)
		})
	}
}

func TestInsert(t *testing.T) {
	sl := initSkipList()

	len := sl.Len()
	sl.Insert(0, "kNew")

	assert.Equal(t, sl.Len(), len+1)
	rank := sl.GetRank(0, "kNew")
	assert.EqualValues(t, rank, 1)
}

func TestGetRank(t *testing.T) {
	sl := initSkipList()

	sl.Insert(1000, "kNew")

	rank := sl.GetRank(1000, "kNew")
	assert.EqualValues(t, rank, sl.Len())
}

func TestUpdateScore(t *testing.T) {
	sl := initSkipList()

	sl.Insert(1000, "kNew")

	assert.EqualValues(t, sl.GetRank(1000, "kNew"), sl.Len())
	assert.EqualValues(t, sl.GetRank(10, "k1"), 1)

	sl.UpdateScore(1000, "kNew", 0)
	assert.EqualValues(t, sl.GetRank(0, "kNew"), 1)
	assert.EqualValues(t, sl.GetRank(10, "k1"), 2)
}

func TestDelete(t *testing.T) {
	sl := initSkipList()
	assert.EqualValues(t, sl.GetRank(80, "k8"), 8)

	sl.Insert(0, "k0")
	assert.EqualValues(t, sl.GetRank(0, "k0"), 1)
	assert.EqualValues(t, sl.GetRank(10, "k1"), 2)
	assert.EqualValues(t, sl.GetRank(80, "k8"), 9)

	sl.Delete(0, "k0")
	assert.EqualValues(t, sl.GetRank(10, "k1"), 1)
	assert.EqualValues(t, sl.GetRank(80, "k8"), 8)
}

func TestGet(t *testing.T) {
	sl := initSkipList()

	assert.NotNil(t, sl.Get(10, "k1"))
	assert.NotNil(t, sl.Get(40, "k4"))
	assert.NotNil(t, sl.Get(80, "k8"))
	assert.Nil(t, sl.Get(80, "kUnknown"))
	assert.Nil(t, sl.Get(100, "k8"))
}
