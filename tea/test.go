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
