package tea

import (
	"testing"
)

func TestNewSelection(t *testing.T) {
	s := NewSelection(Passing("A"))
	if len(s.nodes) != 1 {
		t.Fatalf("expected 1 node in new selection, saw %d", len(s.nodes))
	}

	l := s.nodes[0]
	if len(l.children) != 0 {
		t.Fatalf("new selection should not have any children, but has %d", len(l.children))
	}

	if len(l.xnodes) != 1 {
		t.Fatalf("expected 1 xnode in lnode, saw %d", len(l.xnodes))
	}
}
