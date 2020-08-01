package tea

func NewSelection(test Test) Selection {
	node := newLNode(test, Selection{})
	return Selection{nodes: []*lnode{node}}
}

// Selection represents a set of nodes in our graph.
type Selection struct {
	nodes []*lnode
}

func (s Selection) Child(test Test) Selection {
	node := newLNode(test, s)
	return Selection{nodes: []*lnode{node}}
}

func (s Selection) And(other Selection) Selection {
	included := make(map[int]bool)

	out := make([]*lnode, 0, len(s.nodes)+len(other.nodes))
	for _, n := range append(s.nodes, other.nodes...) {
		if !included[n.id] {
			out = append(out, n)
			included[n.id] = true
		}
	}

	return Selection{nodes: out}
}

// xnodes represents all xnodes in the selected lnodes
func (s Selection) xnodes() []*xnode {
	xnodes := make([]*xnode, 0, s.countXNodes())
	for _, L := range s.nodes {
		for i, _ := range L.xnodes {
			xnodes = append(xnodes, &L.xnodes[i])
		}
	}
	return xnodes
}

func (s Selection) countXNodes() int {
	total := 0
	for _, child := range s.nodes {
		total += len(child.xnodes)
	}
	return total
}

// xleaves looks at all of the selected xnodes, and for every selected xnode,
// traverses the x graph until we arrive at the set of all leaf nodes that have
// a selected ancestor. If the selection consists of the root node, the xleaves
// are all of the leaves of the x graph.
func (s *Selection) xleaves() []*xnode {
	// honestly think that by definition every xnode in the selection has a
	// non-overlapping set of leaves but thinking about this shit is extremely
	// starting to hurt my brain so I'm going to write this in a way that's
	// maybe very redundant.

	seen := make(map[string]bool)
	var leaves []*xnode
	for _, x := range s.xnodes() {
		for _, leaf := range x.leaves() {
			if seen[leaf.label()] {
				panic("double-counting leaves somehow")
			}
			seen[leaf.label()] = true
			leaves = append(leaves, leaf)
		}
	}
	return leaves
}
