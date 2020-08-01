package tea

func NewSelection(test Test) Selection {
	return Selection{nodes: []*lnode{newLNode(test)}}
}

// Selection represents a set of nodes in our graph.
type Selection struct {
	nodes []*lnode
}

func (s Selection) Child(test Test) Selection {
	node := newLNode(test, s.nodes...)
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
		for _, x := range L.xnodes {
			xnodes = append(xnodes, &x)
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

// func (s Selection) writeXDOT(w io.Writer) {
// 	xnodes := s.xnodes()
//
// 	type xedge [2]string
// 	included := make(map[xedge]bool)
// 	edges := make([]xedge, 0, len(xnodes))
//
// 	for _, X := range xnodes {
//
// 		for p := X; p != nil; p = p.parent {
//
// 		}
// 	}
// }
