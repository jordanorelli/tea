package tea

import "testing"

// Run runs a tree of tests, starting from its root.
func Run(t *testing.T, tree *Tree) {
	plan := tree.plan()
	t.Logf("steps in plan: %d", len(plan))
	for _, start := range plan {
		t.Logf("start test %T: %#v", start.Test, start.Test)
		start.run(t)
	}
}

func New(test Test) *Tree {
	return &Tree{Test: test}
}

type Tree struct {
	Test
	parent   *Tree
	children []*Tree
}

func (t *Tree) Child(test Test) *Tree {
	child := &Tree{
		Test:   test,
		parent: t,
	}
	t.children = append(t.children, child)
	return child
}
