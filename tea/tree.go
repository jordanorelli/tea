package tea

import (
	// "reflect"
	"testing"
)

// Run runs a tree of tests, starting from its root.
func Run(t *testing.T, tree *Tree) {
	t.Run(tree.name, func(t *testing.T) {
		setup(t, tree)

		tree.test.Run(t)

		for _, child := range tree.children {
			if t.Failed() || t.Skipped() {
				skip(t, child)
			} else {
				Run(t, child)
			}
		}
	})
}

func setup(t *testing.T, tree *Tree) {
	if tree.parent != nil {
		setup(t, tree.parent)
		tree.parent.test.Run(t)
	}
}

func skip(t *testing.T, tree *Tree) {
	t.Run(tree.name, func(t *testing.T) {
		for _, child := range tree.children {
			skip(t, child)
		}
		t.Skip("tea skipped: dependency failed")
	})
}

func New(test Test) *Tree {
	return &Tree{
		test: test,
		name: parseName(test),
	}
}

type Tree struct {
	test     Test
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

// func clone(t Test) Test {
// 	T := reflect.TypeOf(t)
// 	switch T.Kind() {
// 	case reflect.Struct:
// 	}
// }

func parseName(test Test) string {
	if s, ok := test.(interface{ String() string }); ok {
		return s.String()
	}
	return "???"
}
