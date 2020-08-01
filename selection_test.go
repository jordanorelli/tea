package tea

import (
	"testing"
)

type selectionTest struct {
	label       string
	selection   Selection
	selLNodes   int // number of lnodes in the selection
	selXNodes   int // number of xnodes in the selection
	reachLNodes int // number of lnodes reachable by the selection
	reachXnodes int // number of xnodes reachable by the selection
}

func (test *selectionTest) Run(t *testing.T) {
	if count := len(test.selection.nodes); count != test.selLNodes {
		t.Errorf("expected %d node in selection, saw %d", test.selLNodes, count)
	}

	if count := test.selection.countXNodes(); count != test.selXNodes {
		t.Errorf("expected %d xnode in lnode, saw %d", test.selXNodes, count)
	}
}

func TestSelections(t *testing.T) {
	tests := []selectionTest{
		{
			label:     "new selection",
			selection: NewSelection(A),
			selLNodes: 1,
			selXNodes: 1,
		},
		{
			label:     "root with one child",
			selection: NewSelection(A).Child(B),
			selLNodes: 1,
			selXNodes: 1,
		},
		{
			label:     "two selected roots",
			selection: NewSelection(A).And(NewSelection(B)),
			selLNodes: 2,
			selXNodes: 2,
		},
	}

	add := func(fn func() selectionTest) { tests = append(tests, fn()) }

	add(func() selectionTest {
		root := NewSelection(A)
		b := root.Child(B)
		return selectionTest{
			label:     "root and child selected",
			selection: root.And(b),
			selLNodes: 2,
			selXNodes: 2,
		}
	})

	add(func() selectionTest {
		root := NewSelection(A)
		b := root.Child(B)
		return selectionTest{
			label:     "an optional test",
			selection: root.And(b).Child(C),
			selLNodes: 1,
			selXNodes: 2,
		}
	})

	add(func() selectionTest {
		root := NewSelection(A)
		b := root.Child(B)
		c := root.Child(C)

		return selectionTest{
			label:     "two children selected",
			selection: b.And(c),
			selLNodes: 2,
			selXNodes: 2,
		}
	})

	add(func() selectionTest {
		root := NewSelection(A)
		b := root.Child(B)
		c := root.Child(C)
		return selectionTest{
			label:     "a diamond test",
			selection: b.And(c).Child(D),
			selLNodes: 1,
			selXNodes: 2,
		}
	})

	add(func() selectionTest {
		root := NewSelection(A)
		b := root.Child(B)
		c := root.Child(C)
		d := b.And(c).Child(D)
		return selectionTest{
			label:     "child of a node having multiple parents",
			selection: d.Child(E),
			selLNodes: 1,
			selXNodes: 2,
		}
	})

	for _, test := range tests {
		t.Run(test.label, test.Run)
	}
}
