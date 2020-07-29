// +build std

package incr

import (
	"testing"
)

func TestOnce(t *testing.T) {
	x := 1

	t.Run("increment", func(t *testing.T) {
		x++
		if x != 2 {
			t.Errorf("expected x to be 2, is %d instead", x)
		}
	})
}

func TestTwice(t *testing.T) {
	x := 1

	t.Run("increment", func(t *testing.T) {
		x++
		if x != 2 {
			t.Errorf("expected x to be 2, is %d instead", x)
		}
	})

	// this fails, because both subtests close over the same x.
	t.Run("increment", func(t *testing.T) {
		x++
		if x != 2 {
			t.Errorf("expected x to be 2, is %d instead", x)
		}
	})
}
