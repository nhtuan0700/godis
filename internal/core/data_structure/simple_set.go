package data_structure

type SimpleSet struct {
	dict map[string]struct{}
}

func NewSimpleSet() *SimpleSet {
	return &SimpleSet{
		dict: make(map[string]struct{}, 0),
	}
}

func (s *SimpleSet) Add(members ...string) int {
	added := 0

	for _, m := range members {
		if _, exist := s.dict[m]; !exist {
			s.dict[m] = struct{}{}
			added++
		}
	}

	return added
}

func (s *SimpleSet) Remove(members ...string) int {
	removed := 0
	for _, m := range members {
		if _, exist := s.dict[m]; exist {
			delete(s.dict, m)
			removed++
		}
	}

	return removed
}

func (s *SimpleSet) IsMember(member string) int {
	if _, exist := s.dict[member]; exist {
		return 1
	}

	return 0
}

func (s *SimpleSet) Members() []string {
	m := make([]string, 0)
	for member := range s.dict {
		m = append(m, member)
	}

	return m
}
