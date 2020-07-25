package tea

import (
	"testing"
)

// Run runs a tree of tests, starting from its root.
func Run(t *testing.T, tree *Tree) {
	t.Run(tree.name, func(t *testing.T) {
		var setup []Test

		for root := tree; root.parent != nil; root = root.parent {
			setup = append(setup, root.parent.Test)
		}

		for i, j := 0, len(setup)-1; i < j; i, j = i+1, j-1 {
			setup[i], setup[j] = setup[j], setup[i]
		}

		for _, test := range setup {
			test.Run(t)
		}

		tree.Test.Run(t)

		if t.Failed() || t.Skipped() {
			for _, child := range tree.children {
				skip(t, child)
			}
			return
		}

		for _, child := range tree.children {
			Run(t, child)
		}
	})
}

func skip(t *testing.T, tree *Tree) {
	t.Run(tree.name, func(t *testing.T) {
		t.Skip("tea skipped: dependency failed")
		for _, child := range tree.children {
			skip(t, child)
		}
	})
}

func New(test Test) *Tree {
	return &Tree{
		Test: test,
		name: parseName(test),
	}
}

type Tree struct {
	Test
	name     string
	parent   *Tree
	children []*Tree
}

func (t *Tree) Child(test Test) *Tree {
	child := New(test)
	child.parent = t
	t.children = append(t.children, child)
	return child
}

func parseName(test Test) string {
	if s, ok := test.(interface{ String() string }); ok {
		return s.String()
	}
	return "???"
}
