package tea

import (
	"reflect"
	"strings"
	"testing"
)

// Run runs a tree of tests, starting from its root.
func Run(t *testing.T, tree *Tree) {
	t.Run(tree.name, func(t *testing.T) {
		test := setup(t, tree)
		test.Run(t)

		for _, child := range tree.children {
			if t.Failed() || t.Skipped() {
				skip(t, child)
			} else {
				Run(t, child)
			}
		}
	})
}

func setup(t *testing.T, tree *Tree) Test {
	test := clone(tree.test)
	if tree.parent != nil {
		p := setup(t, tree.parent)
		p.Run(t)
		test = merge(test, p)
	}
	return test
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
	parts := strings.Split(f.Tag.Get("tea"), ",")
	for _, part := range parts {
		if part == "load" {
			return true
		}
	}
	return false
}

func parseName(test Test) string {
	if s, ok := test.(interface{ String() string }); ok {
		return s.String()
	}
	return "???"
}
