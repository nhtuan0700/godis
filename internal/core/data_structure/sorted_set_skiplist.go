package data_structure

type ZSet struct {
	zskiplist *Skiplist
	dict      map[string]float64
}

func NewZSet() *ZSet {
	return &ZSet{
		zskiplist: CreateSkiplist(),
		dict:      make(map[string]float64),
	}
}

func (zs *ZSet) Add(score float64, elm string) bool {
	if len(elm) == 0 {
		return false
	}
	
	_, exist := zs.dict[elm]
	if exist {
		return false
	}

	newNode := zs.zskiplist.Insert(score, elm)
	if newNode == nil {
		return false
	}
	zs.dict[elm] = newNode.score
	return true
}

/*
Return 0-based rank of the object or -1 if the object does not exist.
If reverse is false, rank is computed considering as first element the one
with the lowest score. Otherwise, rank is computed considering as element with rank 0 the
one with the highest score
*/
func (zs *ZSet) GetRank(elm string, reverse bool) (uint64, bool) {
	setSize := len(zs.dict)
	score, exist := zs.dict[elm]
	if !exist {
		return 0, false
	}

	rank := zs.zskiplist.GetRank(score, elm)
	if reverse {
		rank = uint64(setSize) - rank
	} else {
		rank--
	}

	return rank, true
}

func (zs *ZSet) GetScore(elm string) (float64, bool) {
	score, exist := zs.dict[elm]
	return score, exist
}

func (zs *ZSet) Remove(elm string) bool {
	score, exist := zs.dict[elm]
	if !exist {
		return false
	}

	return zs.zskiplist.Delete(score, elm)
}
