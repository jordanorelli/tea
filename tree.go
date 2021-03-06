package tea

import (
	"reflect"
	"strings"
	"testing"
)

// Run runs a tree of tests. Tests will be run recursively starting at the
// provided node and descending to all of its children. All of its parent nodes
// will also be run since they are prerequisites, but none of its sibling node
// will be executed.
//
// Since Run will walk all of the descendents of the provided node, a typical
// usage would be to write a top-level Go test which is a single tree of tea
// Test values. You would then call Run just once, by supplying to Run the root
// node of your tree.
func Run(t *testing.T, tree *Tree) {
	t.Run(tree.name, func(t *testing.T) {
		history, _ := exec(t, tree)
		for _, test := range history {
			if a, ok := test.(After); ok {
				a.After(t)
			}
		}

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

// exec runs the provided test and all of its ancestors in the provided testing
// context. exec returns the environment produced by running these tests.
func exec(t *testing.T, tree *Tree) ([]Test, *env) {
	if tree == nil {
		return nil, nil
	}

	if tree.parent == nil {
		test := clone(tree.test)
		test.Run(t)
		return []Test{test}, mkenv(test)
	}

	history, e := exec(t, tree.parent)
	test := clone(tree.test)
	if err := e.load(test); err != nil {
		t.Errorf("test plan failed: %s", err)
	} else {
		test.Run(t)
	}
	return append([]Test{test}, history...), e.save(test)
}

// skip skips the provided tree node as well as all of its children.
func skip(t *testing.T, tree *Tree) {
	t.Run(tree.name, func(t *testing.T) {
		for _, child := range tree.children {
			skip(t, child)
		}
		t.Skip("tea skipped: dependency failed")
	})
}

// New creates a new testing Tree starting with a root test. Given this root
// Tree node, consumers can add successive nodes to the tree as children of the
// root.
func New(test Test) *Tree {
	return &Tree{
		test: test,
		name: parseName(test),
	}
}

// Tree represents a node in a Tree of tests. Callers create Tree elements in
// one of two ways: by calling New to create a new Tree with the provided test
// as its root, or by calling the Child method on an existing Tree to add a
// child node to the tree.
type Tree struct {
	test     Test
	name     string
	parent   *Tree
	children []*Tree
}

// Child creates a new Tree node as a child of the current tree node, returning
// the newly created child node.
func (t *Tree) Child(test Test) *Tree {
	child := New(test)
	child.parent = t
	t.children = append(t.children, child)
	return child
}

// clone clones a test value, yielding a new test value that can be executed
// and mutated such that the original is not mutated.
func clone(t Test) Test {
	srcV := reflect.ValueOf(t).Elem()
	destV := reflect.New(srcV.Type())
	destV.Elem().Set(srcV)
	return destV.Interface().(Test)
}

// isSaveField takes a struct field and checks its tags for a save tag,
// indicating that the field's value should persist between tests
func isSaveField(f reflect.StructField) bool {
	// PkgPath is empty string when the identifier is unexported.
	if f.PkgPath != "" {
		return false
	}
	parts := strings.Split(f.Tag.Get("tea"), ",")
	for _, part := range parts {
		if part == "save" {
			return true
		}
	}
	return false
}

// isLoadField takes a struct field and checks its tags for a load tag,
// indicating that the field's value should be populated by a saved value from
// a prior test in the chain.
func isLoadField(f reflect.StructField) bool {
	// PkgPath is empty string when the identifier is unexported.
	if f.PkgPath != "" {
		return false
	}
	parts := strings.Split(f.Tag.Get("tea"), ",")
	for _, part := range parts {
		if part == "load" {
			return true
		}
	}
	return false
}

func isMatchField(f reflect.StructField) bool {
	// PkgPath is empty string when the identifier is unexported.
	if f.PkgPath != "" {
		return false
	}
	parts := strings.Split(f.Tag.Get("tea"), ",")
	for _, part := range parts {
		if part == "match" {
			return true
		}
	}
	return false
}

func getMatchFields(t reflect.Type) []reflect.StructField {
	var fields []reflect.StructField
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if isMatchField(f) {
			fields = append(fields, f)
		}
	}
	return fields
}
