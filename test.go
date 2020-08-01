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

// clone clones a test value, yielding a new test value that can be executed
// and mutated such that the original is not mutated. Tests containing pointers
// to objects that were not created by tea will probably not work right. That's
// like, kinda on you though, I can't really enforce things that the Go type
// system doesn't let me enforce.
func clone(t Test) Test {
	v := reflect.ValueOf(t)
	switch v.Kind() {
	case reflect.Ptr:
		v = v.Elem()
	}
	destV := reflect.New(v.Type())
	destV.Elem().Set(v)
	return destV.Interface().(Test)
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

// Pass is a Test value that always passes.
const Pass = Passing("test passed")

// Passing is a Test type that always passes. Every value of the Passing type,
// including the zero value, is a test that will always pass.
type Passing string

func (p Passing) Run(t *testing.T) {}
func (p Passing) String() string   { return string(p) }

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
