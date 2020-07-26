// +build tea

package thing

import (
	"fmt"
	"testing"

	"github.com/jordanorelli/tea"
)

type empty struct{}

func (e *empty) Run(t *testing.T) {}

type testThingSetup struct {
	Thing *Thing `tea:"save"`
}

func (test *testThingSetup) Run(t *testing.T) {
	t.Logf("[%s] running testThingSetup", t.Name())
	if test.Thing != nil {
		t.Fatal("should be nil")
	}
	test.Thing = new(Thing)
}

func (test *testThingSetup) After(t *testing.T) {
	t.Logf("[%s] testThingSetup after", t.Name())
}

func (test testThingSetup) String() string { return "thingSetup" }

type setKey struct {
	Thing *Thing `tea:"load"`
	key   string
	value string
	bad   bool
}

func (test setKey) String() string {
	return fmt.Sprintf("setKey:%s", test.key)
}

func (test *setKey) Run(t *testing.T) {
	t.Logf("[%s] running setKey key: %q value: %q", t.Name(), test.key, test.value)

	// test.Thing is automatically propagated from the prior test by tea!
	err := test.Thing.Set(test.key, test.value)
	if !test.bad && err != nil {
		t.Errorf("should be able to set %q=%q but saw error %v", test.key, test.value, err)
	}
	if test.bad && err == nil {
		t.Errorf("able to set bad values %q=%q", test.key, test.value)
	}
}

func (test *setKey) After(t *testing.T) {
	t.Logf("[%s] setKey after key: %q value: %q", t.Name(), test.key, test.value)
}

func TestThing(t *testing.T) {
	root := tea.New(new(testThingSetup))

	root.Child(&setKey{key: "alice", value: "apple"})
	root.Child(&setKey{key: "bob", value: "banana"})
	root.Child(new(empty)).Child(&setKey{key: "carol", value: "cherry"})
	root.Child(&setKey{bad: true})

	bob := root.Child(&setKey{key: "b ob", value: "banana"})
	bob.Child(&setKey{key: "car-el", value: "cherry"})
	dave := bob.Child(&setKey{key: "dave", value: "durian"})
	dave.Child(&setKey{key: "evan", value: "elderberry"})

	root.Child(&setKey{key: "al ice", bad: true})
	root.Child(&setKey{key: " alice", bad: true})
	root.Child(&setKey{key: "alice ", bad: true})

	tea.Run(t, root)
}
