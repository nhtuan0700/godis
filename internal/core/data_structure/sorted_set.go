package data_structure

type SortedSet struct {
	tree *BPlusTree
}

func CreateZSet(degree int) *SortedSet {
	return &SortedSet{
		tree: NewBPlusTree(degree),
	}
}

func (ss *SortedSet) Add(score float64, member string) int {
	return ss.tree.Add(score, member)  
}

func (ss *SortedSet) GetScore(member string) (float64, bool) {
	return ss.tree.Score(member)
}

func (ss *SortedSet) GetRank(member string) int {
	return ss.tree.GetRank(member)
}
