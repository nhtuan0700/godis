package data_structure_test

import (
	"testing"

	"github.com/nhtuan0700/godis/internal/core/data_structure"
	"github.com/stretchr/testify/assert"
)

var ss *data_structure.SortedSet

func TestMain(m *testing.M) {
	ss = data_structure.CreateZSet(3)
	ss.Add(10.0, "k1")
	ss.Add(20.0, "k2")
	ss.Add(30.0, "k3")
	ss.Add(40.0, "k4")
	ss.Add(50.0, "k5")
	ss.Add(60.0, "k6")
	ss.Add(80.0, "k8")
	ss.Add(70.0, "k7")
	m.Run()
}

func TestSortedSet(t *testing.T) {
	testCases := []struct {
		member string
		score  float64
		rank   int
	}{
		{
			member: "k1",
			score:  10,
			rank:   0,
		},
		{
			member: "k2",
			score:  20,
			rank:   1,
		},
		{
			member: "k3",
			score:  30,
			rank:   2,
		},
		{
			member: "k4",
			score:  40,
			rank:   3,
		},
		{
			member: "k5",
			score:  50,
			rank:   4,
		},
		{
			member: "k6",
			score:  60,
			rank:   5,
		},
		{
			member: "k7",
			score:  70,
			rank:   6,
		},
		{
			member: "k8",
			score:  80,
			rank:   7,
		},
	}

	for _, tc := range testCases {
		t.Run("RankAndScore", func(t *testing.T) {
			score, _ := ss.GetScore(tc.member)
			assert.Equal(t, score, tc.score)
			rank := ss.GetRank(tc.member)
			assert.Equal(t, rank, tc.rank)
		})
	}
}
