package tea

import (
	"testing"
)

func TestSave(t *testing.T) {
	type saveFoo struct {
		empty
		Foo int `tea:"save"`
		Bar string
	}

	type loadFoo struct {
		empty
		Foo int `tea:"load"`
		Bar string
	}

	t.Run("empty begets nil", func(t *testing.T) {
		e := mkenv(new(empty))
		if e != nil {
			t.Errorf("saw unexpected env value looking for nil: %v", e)
		}
	})

	t.Run("unexported fields are ignored", func(t *testing.T) {
		type test struct {
			empty
			foo int `tea:"save"`
		}

		if e := mkenv(test{foo: 5}); e != nil {
			t.Errorf("saw unexpected env value looking for nil: %v", e)
		}
	})

	t.Run("save an int", func(t *testing.T) {
		e := mkenv(&saveFoo{Foo: 5})
		if e == nil {
			t.Fatalf("saw nil env when expecting a valid env")
		}

		if e.key != "Foo" {
			t.Errorf("expected key %q but saw %q instead", "Foo", e.key)
		}

		if e.value != 5 {
			t.Errorf("expected value %v but saw %v instead", 5, e.value)
		}
	})

	t.Run("load an int", func(t *testing.T) {
		e := mkenv(&saveFoo{Foo: 5})
		test := new(loadFoo)

		e.load(test)
		if test.Foo != 5 {
			t.Errorf("expected value %v but saw %v instead", 5, test.Foo)
		}
	})

	t.Run("loads can fail", func(t *testing.T) {
		e := mkenv(new(empty))
		test := new(loadFoo)
		if err := e.load(test); err == nil {
			t.Errorf("expected a load error but did not see one")
		}
	})
}
