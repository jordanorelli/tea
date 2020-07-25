package tea

import "testing"

type step struct {
	Test
	name string
	next *step
	skip bool
}

func (s *step) run(t *testing.T) {
	if s.skip {
		if s.next != nil {
			s.next.skip = true
			t.Run(s.next.name, s.next.run)
		}
		t.SkipNow()
	}

	s.Test.Run(t)
	if s.next != nil {
		if t.Failed() || t.Skipped() {
			s.next.skip = true
		}
		t.Run(s.next.name, s.next.run)
	}
}

func (t *Tree) plan() []step {
	if len(t.children) == 0 {
		s := &step{Test: t.Test, name: parseName(t.Test)}
		for t.parent != nil {
			t = t.parent
			s = &step{Test: t.Test, name: parseName(t.Test), next: s}
		}
		return []step{*s}
	}

	var steps []step
	for _, child := range t.children {
		steps = append(steps, child.plan()...)
	}
	return steps
}
