package tea

import (
	"fmt"
	"reflect"
	"testing"
)

type Test interface {
	Run(*testing.T)
}

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
