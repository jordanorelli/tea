package tea

import (
	"testing"
)

type selectionTest struct {
	label     string
	selection Selection
	lnodes    []string
	xnodes    []string
	xleaves   []string
}

func (test *selectionTest) Run(t *testing.T) {
	LWanted := wantStrings("selected lnode names", test.lnodes...)
	for _, L := range test.selection.nodes {
		LWanted.add(L.name)
	}
	LWanted.report(t)

	XWanted := wantStrings("selected xnode labels", test.xnodes...)
	for _, X := range test.selection.xnodes() {
		XWanted.add(X.label())
	}
	XWanted.report(t)

	// XLeavesWanted := wantStrings("leaf xnode labels", test.xleaves...)
	// for _, X := range test.selection.xnodes() {
	// 	leaves := X.leaves()
	// 	for _, leaf := range leaves {
	// 		XLeavesWanted.

	// 	}

	// }
}

func TestSelections(t *testing.T) {
	tests := []selectionTest{
		{
			label:     "new selection",
			selection: NewSelection(A),
			lnodes:    []string{"A"},
			xnodes:    []string{"A.0"},
		},
		{
			label:     "root with one child",
			selection: NewSelection(A).Child(B),
			lnodes:    []string{"B"},
			xnodes:    []string{"B.0.A"},
		},
		{
			label:     "two selected roots",
			selection: NewSelection(A).And(NewSelection(B)),
			lnodes:    []string{"A", "B"},
			xnodes:    []string{"A.0", "B.0"},
		},
	}

	add := func(fn func() selectionTest) { tests = append(tests, fn()) }

	add(func() selectionTest {
		root := NewSelection(A)
		b := root.Child(B)
		return selectionTest{
			label:     "root and child selected",
			selection: root.And(b),
			lnodes:    []string{"A", "B"},
			xnodes:    []string{"A.0", "B.0.A"},
		}
	})

	add(func() selectionTest {
		root := NewSelection(A)
		b := root.Child(B)
		return selectionTest{
			label:     "an optional test",
			selection: root.And(b).Child(C),
			lnodes:    []string{"C"},
			xnodes:    []string{"C.0.A", "C.1.B"},
		}
	})

	add(func() selectionTest {
		root := NewSelection(A)
		b := root.Child(B)
		c := root.Child(C)

		return selectionTest{
			label:     "two children selected",
			selection: b.And(c),
			lnodes:    []string{"B", "C"},
			xnodes:    []string{"B.0.A", "C.0.A"},
		}
	})

	add(func() selectionTest {
		root := NewSelection(A)
		b := root.Child(B)
		c := root.Child(C)
		return selectionTest{
			label:     "a diamond test",
			selection: b.And(c).Child(D),
			lnodes:    []string{"D"},
			xnodes:    []string{"D.0.B", "D.1.C"},
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
			lnodes:    []string{"E"},
			xnodes:    []string{"E.0.D", "E.1.D"},
		}
	})

	add(func() selectionTest {
		root := NewSelection(A)
		b := root.Child(B)
		c := root.Child(C)
		d := b.And(c).Child(D)
		d.Child(E)
		return selectionTest{
			label:     "the root of a complex graph",
			selection: root,
			lnodes:    []string{"A"},
			xnodes:    []string{"A.0"},
		}
	})

	for _, test := range tests {
		t.Run(test.label, test.Run)
	}
}
