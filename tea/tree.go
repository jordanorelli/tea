package tea

import (
	"testing"
)

// Run runs a tree of tests, starting from its root.
func Run(t *testing.T, tree *Tree) {
	t.Run(tree.name, func(t *testing.T) {
		tree.Test.Run(t)

		if t.Failed() || t.Skipped() {
			for _, child := range tree.children {
				skip(t, child)
			}
			return
		}

		for _, child := range tree.children {
			tree.Test.Run(t)
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
