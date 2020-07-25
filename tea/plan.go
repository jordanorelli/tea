package tea

import "testing"

type step struct {
	Test
	next *step
}

func (s *step) run(t *testing.T) {
	t.Logf("running step: %v", s.Test)
	s.Test.Run(t)
	if s.next != nil {
		s.next.run(t)
	}
}

func (t *Tree) plan() []step {
	if len(t.children) == 0 {
		// this is a leaf node.
		s := &step{Test: t.Test}
		for t.parent != nil {
			t = t.parent
			s = &step{Test: t.Test, next: s}
		}
		return []step{*s}
	}

	var steps []step
	for _, child := range t.children {
		steps = append(steps, child.plan()...)
	}
	return steps
}
