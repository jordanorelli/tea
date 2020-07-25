// +build tea

package thing

import (
	"fmt"
	"testing"

	"./tea"
)

type testThingSetup struct {
	thing *Thing
}

func (test *testThingSetup) Run(t *testing.T) {
	t.Log("Running testThingSetup")
	test.thing = new(Thing)
}

func (test testThingSetup) String() string { return "thingSetup" }

type setKey struct {
	key   string
	value string
	bad   bool
}

func (test setKey) String() string {
	return fmt.Sprintf("setKey(%q=%q)", test.key, test.value)
}

func (test *setKey) Run(t *testing.T) {
	t.Logf("Running setKey key: %q value: %q bad?: %t", test.key, test.value, test.bad)
	thing := new(Thing)

	err := thing.Set(test.key, test.value)
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
	bob := root.Child(&setKey{key: "bob", value: "banana"})
	bob.Child(&setKey{key: "car-el", value: "candy"})
	root.Child(&setKey{key: "d' oh", bad: true})
	tea.Run(t, root)
}
