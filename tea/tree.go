package tea

import (
	// "reflect"
	"testing"
)

// Run runs a tree of tests, starting from its root.
func Run(t *testing.T, tree *Tree) {
	plan := tree.plan()
	for _, step := range plan {
		t.Run(step.name, step.run)
	}
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
	// T := reflect.TypeOf(test)
	// for i := 0; i < T.NumField(); i++ {
	// 	field := T.Field(i)
	// }
}
