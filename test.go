package tea

import (
	"fmt"
	"reflect"
	"testing"
)

// Test is a test value: a value that, given an instance of a testing.T, can be
// used to execute a single test.
type Test interface {
	Run(*testing.T)
}

// After defines the interface used for performing test cleanup. If a Test
// value also implements After, that test's After method will be called after
// all tests are run. Tests in a sequence will have their After methods called
// in the reverse order of their Run methods; a test always runs its After
// method after all of its children have completed their own After methods.
type After interface {
	After(*testing.T)
}

func fail(t string, args ...interface{}) Test {
	return failure{cause: fmt.Errorf(t, args...)}
}

type failure struct {
	cause error
}

func (f failure) Run(t *testing.T) {
	t.Error(f.cause.Error())
}

// empty is an empty test. It does nothing when run, it's just used as a
// sentinel value to create notes in the test graph and for ... testing the tea
// package itself.
type empty struct{}

func (e empty) Run(t *testing.T) {}

// parseName parses the name for a given test
func parseName(test Test) string {
	if s, ok := test.(interface{ String() string }); ok {
		return s.String()
	}

	tv := reflect.ValueOf(test)
	switch tv.Type().Kind() {
	case reflect.Ptr:
		tv = tv.Elem()
	}
	name := tv.Type().Name()
	if name == "" {
		return "unknown-test"
	}
	return name
}
