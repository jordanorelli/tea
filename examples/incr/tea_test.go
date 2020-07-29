// +build !convey,!std

package incr

import (
	"github.com/jordanorelli/tea"
	"testing"
)

type testInt struct {
	X int `tea:"save"`
}

func (testInt) Run(t *testing.T) {}

type testIncr struct {
	X      int `tea:"load"`
	expect int
}

func (test *testIncr) Run(t *testing.T) {
	test.X++
	if test.X != test.expect {
		t.Errorf("expected X to be %d, is %d instead", test.expect, test.X)
	}
}

func TestOnce(t *testing.T) {
	root := tea.New(&testInt{X: 1})
	root.Child(&testIncr{expect: 2})
	tea.Run(t, root)
}

func TestTwice(t *testing.T) {
	root := tea.New(&testInt{X: 1})
	root.Child(&testIncr{expect: 2})
	root.Child(&testIncr{expect: 2})
	tea.Run(t, root)
}
