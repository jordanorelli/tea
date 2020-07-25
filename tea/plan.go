package tea

import "testing"

type step struct {
	Test
	name string
	next *step
	skip bool
}

func (s *step) run(t *testing.T) {
	name := ""
	for s := s; s != nil; s = s.next {
		name += s.name
		if s.next != nil {
			name += "/"
		}
	}
	t.Run(name, func(t *testing.T) {
		for s := s; s != nil; s = s.next {
			s.Test.Run(t)
		}
	})
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
