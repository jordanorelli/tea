package tea

import (
	"testing"
)

type selectionTest struct {
	label     string
	selection Selection
	lnodes    []string
	xnodes    []string
}

func (test *selectionTest) Run(t *testing.T) {
	lfound := make(map[string]bool)
	for _, L := range test.selection.nodes {
		if len(L.parents) > 0 {
			pnames := make([]string, 0, len(L.parents))
			for _, p := range L.parents {
				pnames = append(pnames, p.name)
			}
			t.Logf("found lnode with label %s having %d parents: %s", L.name, len(L.parents), pnames)
		} else {
			t.Logf("found root lnode with label %s", L.name)
		}
		lfound[L.name] = true
	}

	for _, expected := range test.lnodes {
		if lfound[expected] {
			delete(lfound, expected)
		} else {
			t.Errorf("missing expected lnode with label %s", expected)
		}
	}

	for label, _ := range lfound {
		t.Errorf("found unexpected lnode with label %s", label)
	}

	xfound := make(map[string]bool)
	for _, x := range test.selection.xnodes() {
		if x.parent != nil {
			t.Logf("found xnode with label %s having parent %s", x.label(), x.parent.label())
		} else {
			t.Logf("found root xnode with label %s", x.label())
		}
		xfound[x.label()] = true
	}

	for _, expected := range test.xnodes {
		if xfound[expected] {
			delete(xfound, expected)
		} else {
			t.Errorf("missing expected xnode with label %s", expected)
		}
	}

	for label, _ := range xfound {
		t.Errorf("found unexpected xnode with label %s", label)
	}
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

	for _, test := range tests {
		t.Run(test.label, test.Run)
	}
}
