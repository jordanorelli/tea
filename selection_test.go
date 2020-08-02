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

	XLeavesWanted := wantStrings("leaf xnode labels", test.xleaves...)
	for _, X := range test.selection.xnodes() {
		for _, leaf := range X.leaves() {
			XLeavesWanted.add(leaf.label())
		}
	}
	XLeavesWanted.report(t)
}

func TestSelections(t *testing.T) {
	tests := []selectionTest{
		{
			label:     "new selection",
			selection: NewSelection(A),
			lnodes:    []string{"A"},
			xnodes:    []string{"A.0"},
			xleaves:   []string{"A.0"},
		},
		{
			label:     "root with one child",
			selection: NewSelection(A).Child(B),
			lnodes:    []string{"B"},
			xnodes:    []string{"B.0.A"},
			xleaves:   []string{"B.0.A"},
		},
		{
			label:     "two selected roots",
			selection: NewSelection(A).And(NewSelection(B)),
			lnodes:    []string{"A", "B"},
			xnodes:    []string{"A.0", "B.0"},
			xleaves:   []string{"A.0", "B.0"},
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
			xleaves:   []string{"B.0.A"},
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
			xleaves:   []string{"C.0.A", "C.1.B"},
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
			xleaves:   []string{"B.0.A", "C.0.A"},
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
			xleaves:   []string{"D.0.B", "D.1.C"},
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
			xleaves:   []string{"E.0.D", "E.1.D"},
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
			xleaves:   []string{"E.0.D", "E.1.D"},
		}
	})

	//         A
	//        / \
	//       /   \
	//      B     C
	//     / \   / \
	//    /   \ /   \
	//   D     E     F
	//        / \   / \
	//       /   \ /   \
	//      G     H     I
	//      |     |     |
	//      |     |     |
	//      J     K     L
	//      |      \   /
	//      |       \ /
	//      M        N
	//
	add(func() selectionTest {
		root := NewSelection(A)
		b := root.Child(B)
		c := root.Child(C)
		b.Child(D)
		e := b.And(c).Child(E)
		f := c.Child(F)
		e.Child(G).Child(J).Child(M)
		h := e.And(f).Child(H)
		l := f.Child(I).Child(L)
		k := h.Child(K)
		k.And(l).Child(N)
		return selectionTest{
			label:     "criss-crossing",
			selection: root,
			lnodes:    []string{"A"},
			xnodes:    []string{"A.0"},
			xleaves: []string{
				"D.0.B", // A B D
				"M.0.J", // A B E G J M
				"M.1.J", // A C E G J M
				"N.0.K", // A B E H K N
				"N.1.K", // A C E H K N
				"N.2.K", // A C F H K N
				"N.3.L", // A C F I L N
			},
		}
	})

	add(func() selectionTest {
		root := NewSelection(A)
		b := root.Child(B)
		c := root.Child(C)
		b.Child(D)
		e := b.And(c).Child(E)
		f := c.Child(F)
		e.Child(G).Child(J).Child(M)
		h := e.And(f).Child(H)
		l := f.Child(I).Child(L)
		k := h.Child(K)
		k.And(l).Child(N)
		return selectionTest{
			label:     "criss-crossing-partial",
			selection: e,
			lnodes:    []string{"E"},
			xnodes:    []string{"E.0.B", "E.1.C"},
			xleaves: []string{
				"M.0.J", // A B E G J M
				"M.1.J", // A C E G J M
				"N.0.K", // A B E H K N
				"N.1.K", // A C E H K N
			},
		}
	})

	for _, test := range tests {
		t.Run(test.label, test.Run)
	}
}
