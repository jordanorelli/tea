package tea

// this is just a collection of ... reusable assertions for the unit tests for
// tea itself.

import (
	"errors"
	"testing"
)

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

func assertError(t *testing.T, fatal bool, err error, target error) {
	if !errors.Is(err, target) {
		if fatal {
			t.Fatalf("expected error to be %s, instead found: %s", target, err)
		} else {
			t.Errorf("expected error to be %s, instead found: %s", target, err)
		}
	}
}

func assertNoError(t *testing.T, fatal bool, err error) {
	if err != nil {
		if fatal {
			t.Fatalf("encountered unexpected error: %v", err)
		} else {
			t.Fatalf("encountered unexpected error: %v", err)
		}
	}
}
