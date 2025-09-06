package data_structure

type Item struct {
	Score  float64
	Member string
}

type Node struct {
	Items    []*Item // A list of key-value pair, = M - 1
	Children []*Node // Pointers to child nodes, = M
	Parent   *Node   // Pointer to parent node
	IsLeaf   bool    // True if it's a leaf node
	Next     *Node   // For leaf nodes, a pointer to next leaf in the sequence
}

type BPlusTree struct {
	Degree int // the maximum number of children a node can have, = M
	Root   *Node
}

func NewBPlusTree(degree int) *BPlusTree {
	return &BPlusTree{
		Degree: degree,
		Root:   &Node{IsLeaf: true},
	}
}

// Maximum of children capacity, = M - 1
func (t *BPlusTree) Capcity() int {
	return t.Degree - 1
}

// TC: O(log_m(n) + O(n))
// => TODO: can optimize: O(1)
func (t *BPlusTree) Score(member string) (float64, bool) {
	node := t.Root
	// Traverse to the first leaf node
	// We have to search all leaf nodes since we dont know the score
	for !node.IsLeaf {
		node = node.Children[0]
	}

	// Iterate through all leaf nodes using the 'Next' pointer
	for node != nil {
		for _, item := range node.Items {
			if item.Member == member {
				return item.Score, true
			}
		}

		node = node.Next
	}

	// member not found
	return 0, false
}

// find the correct leaf node
// append to right position
func (t *BPlusTree) Add(score float64, member string) int {
	newItem := &Item{Member: member, Score: score}

	if len(member) == 0 {
		return 0
	}

	node := t.Root
	for !node.IsLeaf {
		i := 0
		// Find the correct child based on the scores
		// TODO: binary search
		for i < len(node.Items) && score >= node.Items[i].Score {
			i++
		}
		node = node.Children[i]
	}

	// check if the member already existed in the leaf node, only update score
	for _, item := range node.Items {
		if item.Member == member {
			item.Score = score
			return 1
		}
	}

	//  Member does node exist, insert it into the sorted position
	insertIdx := 0
	// TODO: binary search
	for insertIdx < len(node.Items) && score >= node.Items[insertIdx].Score {
		insertIdx++
	}
	node.Items = append(node.Items[:insertIdx], append([]*Item{newItem}, node.Items[insertIdx:]...)...)

	// split the node if it's over capacity
	if len(node.Items) > t.Capcity() {
		t.splitNode(node)
	}
	return 1
}

func (t *BPlusTree) splitNode(node *Node) {
	// If the node is the root, we need to create a new root.
	if node.Parent == nil {
		t.splitRoot()
		return
	}

	// Split based on whether the node is a leaf or an internal node.
	if node.IsLeaf {
		t.splitLeaf(node)
	} else {
		t.splitInternal(node)
	}
}

func (t *BPlusTree) splitLeaf(node *Node) {
	// create a new sibling leaf node
	newLeaf := &Node{
		IsLeaf: true,
		Parent: node.Parent,
		Next:   node.Next,
	}

	medianIdx := len(node.Items) / 2
	// move the second half of the items to the new leaf
	newLeaf.Items = append(newLeaf.Items, node.Items[medianIdx:]...)
	node.Items = node.Items[:medianIdx]
	// update the 'Next' pointer for sequential traversal
	node.Next = newLeaf

	// promote the first key of the new leaf to the parent
	parent := node.Parent
	promotedItem := newLeaf.Items[0]

	// find the insertion point in the parent nodes
	childIdx := 0
	for childIdx < len(parent.Children) && parent.Children[childIdx] != node {
		childIdx++
	}

	// insert promoted key and the new child node into the parent
	parent.Items = append(parent.Items[:childIdx], append([]*Item{promotedItem}, parent.Items[childIdx:]...)...)
	parent.Children = append(parent.Children[:childIdx+1], append([]*Node{newLeaf}, parent.Children[childIdx+1:]...)...)

	// if parent now overflows, split it too.
	if len(parent.Items) > t.Capcity() {
		t.splitNode(parent)
	}
}

func (t *BPlusTree) splitInternal(node *Node) {
	// Create a new sibling internal node
	newInternal := &Node{
		Parent: node.Parent,
		IsLeaf: false,
	}

	medianIdx := len(node.Items) / 2
	// Protomoted the median key to thge parent
	promotedItem := node.Items[medianIdx]

	// Move the second half of the items and children to the new node
	newInternal.Items = append(newInternal.Items, node.Items[medianIdx+1:]...)
	newInternal.Children = append(newInternal.Children, node.Children[medianIdx+1:]...)

	// Trim original node
	node.Items = node.Items[:medianIdx]
	node.Children = node.Children[:medianIdx+1]

	// Update parent pointers for the new children.
	for _, child := range newInternal.Children {
		child.Parent = newInternal
	}

	// Now insert the promoted key and the new node into the parent
	parent := node.Parent
	childIdx := 0
	for childIdx < len(parent.Items) && parent.Children[childIdx] != node {
		childIdx++
	}

	// Insert the promoted key and the new child
	parent.Items = append(parent.Items[:childIdx], append([]*Item{promotedItem}, parent.Items[childIdx:]...)...)
	parent.Children = append(parent.Children[:childIdx+1], append([]*Node{newInternal}, parent.Children[childIdx+1:]...)...)

	if len(parent.Items) > t.Capcity() {
		t.splitNode(parent)
	}
}

func (t *BPlusTree) splitRoot() {
	oldRoot := t.Root
	newRoot := &Node{}

	t.Root = newRoot
	// Set the old root as the first child of the new root
	oldRoot.Parent = newRoot
	newRoot.Children = append(newRoot.Children, oldRoot)

	if oldRoot.IsLeaf {
		t.splitLeaf(oldRoot)
	} else {
		t.splitInternal(oldRoot)
	}
}

func (t *BPlusTree) GetRank(member string) int {
	rank := 0

	// Find the first leaf node
	node := t.Root
	for !node.IsLeaf {
		node = node.Children[0] // always go to the left most child
	}

	// Traversal all leaf nodes from the beginning
	for node != nil {
		for _, item := range node.Items {
			// check if we have found the member
			if item.Member == member {
				return rank // Return the current rank
			}
			rank++
		}

		node = node.Next
	}

	return -1 // member not found
}
