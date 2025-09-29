package data_structure_test

import (
	"testing"

	"github.com/nhtuan0700/godis/internal/core/data_structure"
	"github.com/stretchr/testify/assert"
)

func TestSortedSetSkiplist(t *testing.T) {
	zs := data_structure.NewZSet()
	zs.Add(10.0, "k1")
	zs.Add(20.0, "k2")
	zs.Add(30.0, "k3")
	zs.Add(40.0, "k4")
	zs.Add(50.0, "k5")
	zs.Add(60.0, "k6")
	zs.Add(80.0, "k8")
	zs.Add(70.0, "k7")
	testCases := []struct {
		member string
		score  float64
		rank   uint64
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
			rank, _ := zs.GetRank(tc.member, false)
			score, _ := zs.GetScore(tc.member)
			assert.Equal(t, score, tc.score)
			assert.Equal(t, rank, tc.rank)
		})
	}
}
