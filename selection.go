package tea

import (
	"testing"
)

func NewSelection(test Test) Selection {
	n := node{test: test}
	return Selection{
		nodes: []*node{&n},
	}
}

func RunSelection(t *testing.T, s Selection) {

}

type Selection struct {
	nodes []*node
}

func (s Selection) Child(test Test) Selection {
	child := &node{
		id:      nextNodeID(),
		test:    test,
		name:    parseName(test),
		parents: s.nodes,
	}
	for _, sn := range s.nodes {
		sn.children = append(sn.children, child)
	}
	return Selection{nodes: []*node{child}}
}

func (s Selection) And(other Selection) Selection {
	included := make(map[int]bool)

	out := make([]*node, 0, len(s.nodes)+len(other.nodes))
	for _, n := range append(s.nodes, other.nodes...) {
		if !included[n.id] {
			out = append(out, n)
			included[n.id] = true
		}
	}

	return Selection{nodes: out}
}
