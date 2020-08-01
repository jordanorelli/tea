package tea

import "testing"

type labelChecker struct {
	name   string
	wanted map[string]bool // all the strings we want
	found  map[string]bool // all the strings we've found
}

func wantStrings(name string, want ...string) labelChecker {
	l := newLabelChecker(name)
	l.want(want...)
	return l
}

func newLabelChecker(name string) labelChecker {
	return labelChecker{
		name:   name,
		wanted: make(map[string]bool),
		found:  make(map[string]bool),
	}
}

func (l *labelChecker) want(names ...string) {
	for _, name := range names {
		l.wanted[name] = true
	}
}

func (l *labelChecker) add(name string) {
	l.found[name] = true
}

func (l *labelChecker) report(t *testing.T) {
	for name, _ := range l.found {
		if l.wanted[name] {
			t.Logf("%s saw expected value %s", l.name, name)
		} else {
			t.Errorf("%s saw unexpected value %s", l.name, name)
		}
	}
	for name, _ := range l.wanted {
		if !l.found[name] {
			t.Errorf("%s missing expected value %s", l.name, name)
		}
	}
}
