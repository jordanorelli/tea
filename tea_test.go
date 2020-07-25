// +build tea

package thing

import (
	"fmt"
	"testing"

	"./tea"
)

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

	err := test.Thing.Set(test.key, test.value)
	if !test.bad && err != nil {
		t.Errorf("should be able to set %q=%q but saw error %v", test.key, test.value, err)
	}
	if test.bad && err == nil {
		t.Errorf("able to set bad values %q=%q", test.key, test.value)
	}
}

func TestThing(t *testing.T) {
	root := tea.New(new(testThingSetup))

	root.Child(&setKey{key: "alice", value: "apple"})
	root.Child(&setKey{key: "bob", value: "banana"})
	root.Child(&setKey{key: "carol", value: "cherry"})

	bob := root.Child(&setKey{key: "b ob", value: "banana"})
	bob.Child(&setKey{key: "car-el", value: "cherry"})
	dave := bob.Child(&setKey{key: "dave", value: "durian"})
	dave.Child(&setKey{key: "evan", value: "elderberry"})

	root.Child(&setKey{key: "al ice", bad: true})
	root.Child(&setKey{key: " alice", bad: true})
	root.Child(&setKey{key: "alice ", bad: true})

	tea.Run(t, root)
}
