package tea

import (
	"io"
)

func NewSelection(test Test) Selection {
	x := xnode{
		test: clone(test),
	}
	l := lnode{
		id:     nextNodeID(),
		name:   parseName(test),
		xnodes: []*xnode{&x},
	}
	x.lnode = &l

	return Selection{
		nodes: []*lnode{&l},
	}
}

// Selection represents a set of nodes in our graph.
type Selection struct {
	nodes []*lnode
}

func (s Selection) Child(test Test) Selection {
	child := &lnode{id: nextNodeID(), name: parseName(test)}
	for _, l := range s.nodes {
		l.child(child, test)
	}
	return Selection{nodes: []*lnode{child}}
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

func (s Selection) countXNodes() int {
	total := 0
	for _, child := range s.nodes {
		total += len(child.xnodes)
	}
	return total
}

func (s Selection) writeDOT(w io.Writer) {

}
