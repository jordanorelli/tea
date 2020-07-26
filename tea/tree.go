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
func Run(t *testing.T, tree *Tree) {
	t.Run(tree.name, func(t *testing.T) {
		setup(t, tree)

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

// setup runs all of the tests ancestor to the given tree, building up a
// testing environment from their side-effects
func setup(t *testing.T, tree *Tree) *env {
	if tree == nil {
		return nil
	}

	if tree.parent == nil {
		test := clone(tree.test)
		test.Run(t)
		return mkenv(test)
	}

	e := setup(t, tree.parent)
	test := clone(tree.test)
	e.load(test)
	test.Run(t)
	return e.save(test)
}

// setup runs all of the dependencies for a given test. All of the tests are
// run in the same subtest (and therefore same goroutine).
// func setup(t *testing.T, tree *Tree) Test {
// 	// clone the user's values before doing anything, we don't want to pollute
// 	// the planning tree.
// 	test := clone(tree.test)
//
// 	if tree.parent != nil {
// 		p := setup(t, tree.parent)
// 		p.Run(t)
// 		test = merge(test, p)
// 	}
//
// 	return test
// }

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

// clone clones a test value, yielding a new test value that can be executed
// and mutated such that the original is not mutated.
func clone(t Test) Test {
	srcV := reflect.ValueOf(t).Elem()
	destV := reflect.New(srcV.Type())
	destV.Elem().Set(srcV)
	return destV.Interface().(Test)
}

// merge merges into dest the fields on the src test that are marked as worth
// saving
func merge(dest Test, src Test) Test {
	destV := reflect.ValueOf(dest).Elem()
	srcV := reflect.ValueOf(src).Elem()

	for i := 0; i < srcV.NumField(); i++ {
		sf := srcV.Type().Field(i)
		if isSaveField(sf) {
			df, ok := destV.Type().FieldByName(sf.Name)
			if ok && isLoadField(df) {
				if sf.Type == df.Type {
					sfv := srcV.FieldByName(sf.Name)
					dfv := destV.FieldByName(sf.Name)
					dfv.Set(sfv)
				}
			}
		}
	}
	return dest
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

// parseName parses the name for a given test
func parseName(test Test) string {
	if s, ok := test.(interface{ String() string }); ok {
		return s.String()
	}
	return "???"
}
