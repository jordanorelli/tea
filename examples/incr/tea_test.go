// +build !convey,!std

package incr

import (
	"github.com/jordanorelli/tea"
	"testing"
)

// testInt is a Test that does nothing. It always passes. It acts as a
// container for a value to be made available to future tests.
type testInt struct {
	// the "save" struct tag instructs tea to save this field for future tests.
	X int `tea:"save"`
}

// Run satisfies the tea.Test interface so that testInt may be used as a
// tea.Test
func (test *testInt) Run(t *testing.T) {
	t.Logf("saving to future tests X = %d", test.X)
}

// testIncr increments an integer and checks that the result of incrementing
// the integer is an expected value.
type testIncr struct {
	// the "load" struct tag instructs tea to load the value of this field from
	// previous tests in this run. Like before, we also use a "save" tag.
	X      int `tea:"load,save"`
	expect int
}

// Run satisfies the tea.Test interface
func (test *testIncr) Run(t *testing.T) {
	t.Logf("loaded from parent tests X = %d", test.X)
	test.X++
	if test.X != test.expect {
		t.Errorf("expected X to be %d, is %d instead", test.expect, test.X)
	}
	t.Logf("saving to future tests X = %d", test.X)
}

func TestOnce(t *testing.T) {
	// we use testInt with X set to 1 as our starting test.
	root := tea.New(&testInt{X: 1})

	// after that test passes, we want to run a test of type testIncr.
	root.Child(&testIncr{expect: 2})

	tea.Run(t, root)
}

func TestTwice(t *testing.T) {
	// same setup as in TestOnce
	root := tea.New(&testInt{X: 1})

	// just like before, we want to run a testIncr test after the root test
	// passes.
	root.Child(&testIncr{expect: 2})

	// This testIncr and the other testIncr are siblings, since we called Child
	// from the same node.
	root.Child(&testIncr{expect: 2})

	tea.Run(t, root)
}

func TestTwiceSeries(t *testing.T) {
	// we use testInt with X set to 1 as our starting test.
	root := tea.New(&testInt{X: 1})

	// by now this should look familiar. The difference here is that we save
	// the newly created node in the tree.
	two := root.Child(&testIncr{expect: 2})

	// adding a child to two here
	two.Child(&testIncr{expect: 3})

	// but of course, we can still add new children to the root.
	root.Child(&testIncr{expect: 2})

	tea.Run(t, root)
}
