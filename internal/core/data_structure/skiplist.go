package data_structure

import (
	"math"
	"math/rand"
	"strings"
)

const (
	SkiplistMaxLevel = 32
)

type SkiplistNode struct {
	elm      string
	score    float64
	backward *SkiplistNode
	levels   []SkiplistLevel
}

type SkiplistLevel struct {
	forward *SkiplistNode
	// span is the number of nodes between current node the node->forward at current level
	span uint64
}

type Skiplist struct {
	head   *SkiplistNode
	tail   *SkiplistNode
	length uint64
	level  int
}

func (sl *Skiplist) CreateNode(level int, score float64, elm string) *SkiplistNode {
	return &SkiplistNode{
		elm:      elm,
		score:    score,
		backward: nil,
		levels:   make([]SkiplistLevel, level),
	}
}

func CreateSkiplist() *Skiplist {
	sl := &Skiplist{
		level:  1,
		length: 0,
	}
	sl.head = sl.CreateNode(SkiplistMaxLevel, math.Inf(-1), "")
	sl.head.backward = nil
	sl.tail = nil
	return sl
}

func (sl *Skiplist) randomLevel() int {
	level := 1
	for rand.Intn(2) == 1 {
		level++
	}

	if level > SkiplistMaxLevel {
		return SkiplistMaxLevel
	}
	return level
}

func (sl *Skiplist) Len() uint64 {
	// exclude the head node
	return sl.length
}

// Insert insert a new element to the SkipList, we allow duplidated scores
// Insert adds a new node with a given score and element to the skiplist.
// The new node is determined probabilistically
func (sl *Skiplist) Insert(score float64, elm string) *SkiplistNode {
	// `update` stores the nodes that need to have their 'forward' pointers updated
	// at each level to insert the new node.
	// `rank` stores the number of nodes visited at each level while searching for
	// the insertion position.
	update := [SkiplistMaxLevel]*SkiplistNode{}
	rank := [SkiplistMaxLevel]uint64{}
	x := sl.head

	// traverse from the highest level down to find the insertion point.
	// this loop determines the `update` and `rank` array
	for i := sl.level - 1; i >= 0; i-- {
		// Initalize rank for the current level based on the previous level's rank
		if i == sl.level-1 {
			rank[i] = 0
		} else {
			rank[i] = rank[i+1]
		}

		// Move forward on the current level as soon as the next node's score is less than
		// or equal to the new node's score.
		// `strings.Compare` handles the case of equals score to maintain a stable sort order.
		for x.levels[i].forward != nil && (x.levels[i].forward.score < score ||
			(x.levels[i].forward.score == score && strings.Compare(x.levels[i].forward.elm, elm) == -1)) {
			// Accumulate the 'span' of each node to calculate the rank
			rank[i] += x.levels[i].span
			// Move to the next node.
			x = x.levels[i].forward
		}

		// Store the last node visited at this level before dropping down.
		update[i] = x
	}

	// Determine the level of the new node using a probabilistic method.
	level := sl.randomLevel()
	if level > sl.level {
		for i := sl.level; i < level; i++ {
			rank[i] = 0
			update[i] = sl.head
			// the span for new levels from the head to the end is the entire list length.
			update[i].levels[i].span = sl.length
		}

		// Update the skiplist's high level
		sl.level = level
	}

	// Create new node.
	x = sl.CreateNode(level, score, elm)
	// Link the new node into the skiplist at all its levels.
	for i := 0; i < level; i++ {
		// Update the forward pointers to insert the new node
		x.levels[i].forward = update[i].levels[i].forward
		update[i].levels[i].forward = x
		// Calculate the span for new node
		x.levels[i].span = update[i].levels[i].span - (rank[0] - rank[i])
		// Update the span for the node before the new node.
		update[i].levels[i].span = rank[0] - rank[i] + 1
	}

	// increase span for untouched level because we have a new node.
	// For levels higher than the new node's level, the span of the `update` nodes
	// (which are the nodes before the insertion point) need to be increased by one.
	// This is because the new node is inserted below them.
	for i := level; i < sl.level; i++ {
		update[i].levels[i].span++
	}

	// Update the backward pointer for the new node, which is at the bottom level(0).
	if update[0] == sl.head {
		x.backward = nil
	} else {
		x.backward = update[0]
	}

	// Update the backward pointer of the node that comes after the new node.
	if x.levels[0].forward != nil {
		x.levels[0].forward.backward = x
	} else {
		// If the new node is the last one in the list, it will become the tail of list
		sl.tail = x
	}

	// Increate total length of the skiplist.
	sl.length++

	// Return the newly inserted node.
	return x
}

/*
Find the rank for an element by both score and key.
Return 0 if the element cannot be found, rank otherwise.
Note that the rank is 1-based because the zsl->head is the first element.
*/
func (sl *Skiplist) GetRank(score float64, elm string) uint64 {
	x := sl.head
	var rank uint64 = 0
	// Traverse from highest level down the bottom level to find the element
	// Calculate the rank at each level that we're traversing

	for i := sl.level - 1; i >= 0; i-- {
		// Move on the current level to the next node as soon as the next node's score is less than
		// or equal to the score.
		for x.levels[i].forward != nil && (x.levels[i].forward.score < score ||
			(x.levels[i].forward.score == score && strings.Compare(x.levels[i].forward.elm, elm) <= 0)) {
			rank += x.levels[i].span
			x = x.levels[i].forward
		}
		if x.score == score && x.elm == elm {
			return rank
		}
	}

	return 0
}

// Update score will update the node with curScore and elm
func (sl *Skiplist) UpdateScore(curScore float64, elm string, newScore float64) *SkiplistNode {
	x := sl.head
	update := [SkiplistMaxLevel]*SkiplistNode{}

	for i := sl.level - 1; i >= 0; i-- {
		for x.levels[i].forward != nil && (x.levels[i].forward.score < curScore ||
			(x.levels[i].forward.score == curScore && strings.Compare(x.levels[i].forward.elm, elm) == -1)) {
			x = x.levels[i].forward
		}

		update[i] = x
	}

	x = x.levels[0].forward
	if x == nil || x.score != curScore || x.elm != elm {
		return nil
	}
	if (x.backward == nil || x.backward.score < newScore) &&
		(x.levels[0].forward == nil || x.levels[0].forward.score > newScore) {
		x.score = newScore
		return x
	}

	sl.DeleteNode(x, update)
	return sl.Insert(newScore, elm)
}

func (sl *Skiplist) DeleteNode(x *SkiplistNode, update [SkiplistMaxLevel]*SkiplistNode) {
	for i := sl.level - 1; i >= 0; i-- {
		// descrease span at each level for the nodes (not deleting node)
		if update[i].levels[i].forward == x {
			update[i].levels[i].forward = x.levels[i].forward
			update[i].levels[i].span += x.levels[i].span - 1
		} else {
			update[i].levels[i].span--
		}
	}

	// update backward for the forward node of deleting node
	if x.levels[0].forward != nil {
		x.levels[0].forward.backward = x.backward
	} else {
		// x is tail
		sl.tail = x.backward
	}

	// decrease the levels if levels contains the only node
	for sl.level > 1 && sl.head.levels[sl.level-1].forward == nil {
		sl.level--
	}

	sl.length--
}

func (sl *Skiplist) Delete(score float64, elm string) bool {
	x := sl.head
	// `update` store the nodes that need to have their forward pointers updated
	// at each level to delete the node
	update := [SkiplistMaxLevel]*SkiplistNode{}

	for i := sl.level - 1; i >= 0; i-- {
		for x.levels[i].forward != nil && (x.levels[i].forward.score < score ||
			(x.levels[i].forward.score == score && strings.Compare(x.levels[i].forward.elm, elm) == -1)) {
			x = x.levels[i].forward
		}

		update[i] = x
	}

	x = x.levels[0].forward
	if x == nil || x.score != score || x.elm != elm {
		return false
	}

	sl.DeleteNode(x, update)
	return true
}

func (sl *Skiplist) Get(score float64, elm string) *SkiplistNode {
	x := sl.head

	for i := sl.level - 1; i > 0; i-- {
		for x.levels[i].forward != nil && (x.levels[i].forward.score < score ||
			(x.levels[i].forward.score == score && strings.Compare(x.levels[i].forward.elm, elm) <= 0)) {
				x = x.levels[i].forward
		}

		if x != nil && x.score == score && x.elm == elm {
			return x
		}
	}

	return nil
}
