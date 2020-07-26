package tea

import (
	"fmt"
	"testing"
)

type Test interface {
	Run(*testing.T)
}

func fail(t string, args ...interface{}) Test {
	return failure{cause: fmt.Errorf(t, args...)}
}

type failure struct {
	cause error
}

func (f failure) Run(t *testing.T) { t.Error(f.cause.Error()) }

// empty is an empty test. It does nothing when run, it's just used as a
// sentinel value to create notes in the test graph and for ... testing the tea
// package itself.
type empty struct{}

func (e empty) Run(t *testing.T) {}
