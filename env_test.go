package tea

import (
	"testing"
)

func TestSave(t *testing.T) {
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

	t.Run("create an env from a test", func(t *testing.T) {
		test := struct {
			empty
			Foo int `tea:"save"`
		}{
			Foo: 5,
		}

		e := mkenv(&test)
		if e == nil {
			t.Fatalf("saw nil env when expecting a valid env")
		}

		foo, ok := e.data["Foo"]
		if !ok {
			t.Errorf("expected field Foo to be saved but was not saved")
		}

		if foo != 5 {
			t.Errorf("expected value %v but saw %v instead", 5, foo)
		}
	})

	t.Run("update an existing env", func(t *testing.T) {
		test := struct {
			empty
			Foo int `tea:"save"`
		}{
			Foo: 5,
		}

		e := mkenv(&test)
		if e == nil {
			t.Fatalf("saw nil env when expecting a valid env")
		}

		foo, ok := e.data["Foo"]
		if !ok {
			t.Errorf("expected field Foo to be saved but was not saved")
		}

		if foo != 5 {
			t.Errorf("expected value %v but saw %v instead", 5, foo)
		}
	})
}

func TestLoad(t *testing.T) {
	t.Run("load an int", func(t *testing.T) {
		e := &env{
			data: map[string]interface{}{"Foo": 5},
		}

		var test struct {
			empty
			Foo int `tea:"load"`
		}

		e.load(&test)
		if test.Foo != 5 {
			t.Errorf("expected value %v but saw %v instead", 5, test.Foo)
		}
	})

	t.Run("loads can fail", func(t *testing.T) {
		e := &env{
			data: map[string]interface{}{"NotFoo": 5},
		}

		var test struct {
			empty
			Foo int `tea:"load"`
		}

		if err := e.load(&test); err == nil {
			t.Errorf("expected a load error but did not see one")
		}
	})
}

func TestMatch(t *testing.T) {
	t.Run("required match field not present", func(t *testing.T) {
		e := &env{
			data: map[string]interface{}{"Foo": 5},
		}

		var test struct {
			empty
			Name string `tea:"match"`
			Foo  int    `tea:"load"`
		}

		if err := e.load(&test); err == nil {
			t.Errorf("expected a load error but did not see one")
		}
	})

	t.Run("required match field has wrong value", func(t *testing.T) {
		e := &env{
			data: map[string]interface{}{
				"Foo":  5,
				"Name": "alice",
			},
		}

		var test struct {
			empty
			Name string `tea:"match"`
			Foo  int    `tea:"load"`
		}
		test.Name = "bob"

		if err := e.load(&test); err == nil {
			t.Errorf("expected a load error but did not see one")
		}
	})

	t.Run("simple match", func(t *testing.T) {
		e := &env{
			data: map[string]interface{}{
				"Foo":  5,
				"Name": "alice",
			},
		}

		var test struct {
			empty
			Name string `tea:"match"`
			Foo  int    `tea:"load"`
		}
		test.Name = "alice"

		if err := e.load(&test); err != nil {
			t.Errorf("unexpected load error: %v", err)
		}
		if test.Foo != 5 {
			t.Errorf("expected Foo to load 5 but is %d instead", test.Foo)
		}
	})

	t.Run("ancestor match", func(t *testing.T) {
		e := &env{
			data: map[string]interface{}{
				"Foo":  3,
				"Name": "bob",
			},
			parent: &env{
				data: map[string]interface{}{
					"Foo":  5,
					"Name": "alice",
				},
			},
		}

		var test struct {
			empty
			Name string `tea:"match"`
			Foo  int    `tea:"load"`
		}
		test.Name = "alice"

		if err := e.load(&test); err != nil {
			t.Errorf("unexpected load error: %v", err)
		}
		if test.Foo != 5 {
			t.Errorf("expected Foo to load 5 but is %d instead", test.Foo)
		}
	})
}

// A.Optional(B).Child(C)
//
//       A           A
//      /|          / \
//     / |         /   \
//    B  | ---->  B     C'
//     \ |        |
//      \|        |
//       C        C

// what to call this thing?
//
//       A               A
//      / \             / \
//     /   \           /   \
//    B     C  ---->  B     C
//     \   /          |     |
//      \ /           |     |
//       D            D     D'
