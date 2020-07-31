package tea

import (
	"testing"
)

func TestSave(t *testing.T) {
	t.Run("empty begets nil", func(t *testing.T) {
		e := mkenv(Pass)
		if e != nil {
			t.Errorf("saw unexpected env value looking for nil: %v", e)
		}
	})

	t.Run("unexported fields are ignored", func(t *testing.T) {
		type test struct {
			Passing
			foo int `tea:"save"`
		}

		if e := mkenv(test{foo: 5}); e != nil {
			t.Errorf("saw unexpected env value looking for nil: %v", e)
		}
	})

	t.Run("create an env from a test", func(t *testing.T) {
		test := struct {
			Passing
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
			Passing
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
			Passing
			Foo int `tea:"load"`
		}

		if err := e.load(&test); err != nil {
			t.Errorf("unexpected load error: %v", err)
		}
		if test.Foo != 5 {
			t.Errorf("expected value %v but saw %v instead", 5, test.Foo)
		}
	})

	t.Run("loads can fail", func(t *testing.T) {
		e := &env{
			data: map[string]interface{}{"NotFoo": 5},
		}

		var test struct {
			Passing
			Foo int `tea:"load"`
		}

		if err := e.load(&test); err == nil {
			t.Errorf("expected a load error but did not see one")
		}
	})

	t.Run("skip load if field is already set", func(t *testing.T) {
		e := &env{
			data: map[string]interface{}{"Foo": 3},
		}

		var test struct {
			Passing
			Foo int `tea:"load"`
		}
		test.Foo = 5

		if err := e.load(&test); err != nil {
			t.Errorf("unexpected load error: %v", err)
		}
		if test.Foo != 5 {
			t.Errorf("load overwrote expected value of 5 with %d", test.Foo)
		}
	})
}

func TestMatch(t *testing.T) {
	t.Run("required match field not present", func(t *testing.T) {
		e := &env{
			data: map[string]interface{}{"Foo": 5},
		}

		var test struct {
			Passing
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
			Passing
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
			Passing
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
			Passing
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

	t.Run("complicated match", func(t *testing.T) {
		type connect struct {
			Passing
			Role string `tea:"save"`
			Name string `tea:"save"`
			ID   int    `tea:"save"`
		}

		type request struct {
			Passing
			Role string `tea:"match"`
			Name string `tea:"match"`
			ID   int    `tea:"load"`
		}

		e := mkenv(connect{
			Role: "host",
			ID:   1,
		})
		e = e.save(connect{
			Role: "player",
			Name: "alice",
			ID:   2,
		})
		e = e.save(connect{
			Role: "player",
			Name: "bob",
			ID:   3,
		})

		bob := request{Role: "player", Name: "bob"}
		if err := e.load(&bob); err != nil {
			t.Errorf("failed to load bob: %s", err)
		} else {
			if bob.ID != 3 {
				t.Errorf("expected bob to have ID 3, has %d instead", bob.ID)
			}
		}

		alice := request{Role: "player", Name: "alice"}
		if err := e.load(&alice); err != nil {
			t.Errorf("failed to load alice: %s", err)
		} else {
			if alice.ID != 2 {
				t.Errorf("expected alice to have ID 2, has %d instead", alice.ID)
			}
		}

		host := request{Role: "host"}
		if err := e.load(&host); err != nil {
			t.Errorf("failed to load host: %s", err)
		} else {
			if host.ID != 1 {
				t.Errorf("expected host to have ID 1, has %d instead", host.ID)
			}
		}

	})
}

// Constructing a test node that has multiple parents:
// -----------------------------------------------------------------------------
//
// In this example, B is an optional test.
//
//    Logical    Execution
//
//       A           A
//      /|          / \
//     / |         /   \
//    B  | ---->  B     C'
//     \ |        |
//      \|        |
//       C        C
//
// This logical graph of test dependencies would yield an execution plan
// consisting of two test chains:
//
//   A -> B -> C
//   A -> C
//
// We could write this as follows:
//
//   root := New(A)
//   b := root.Child(B)
//   root.And(b).Child(C)
//
// Alternatively:
//
//   root := New(A)
//   root.Child(B).And(root).Child(C)
//
// If we permit a selection to append multiple children, we could write this as
// follows:
//
//   root := New(A)
//   root.Child(B, Pass).Child(C)
//
//   This last form is not strictly the same, since it includes an additional
//   node in the graph which is a passing test. However since Pass is a
//   specific example, we can trivially remove nodes having a test value of
//   Pass in the planning phase.
//
// Another simple example: a diamond-shaped test graph
//
//    Logical         Execution
//
//       A               A
//      / \             / \
//     /   \           /   \
//    B     C  ---->  B     C
//     \   /          |     |
//      \ /           |     |
//       D            D     D'
//
// Test Plan:
//   - A -> B -> D
//   - A -> C -> D
//
// Expressed in test code as follows:
//
//   root := New(A)
//   both := root.Child(B, C)
//   both.Child(D)
//
// Alternatively:
//
//   New(A).Child(B, C).Child(D)
//
//
//
// This API is fairly straightforward to use, but breaks down with even simple
// shapes:
//
//         A
//        / \
//       /   \
//      B     C
//     / \   /
//    /   \ /
//   E     D
//
// Test Plan:
//   - A -> B -> E
//   - A -> B -> D
//   - A -> C -> D
//
// Essentially what we're saying is:
//   Run test A.
//   If test A passes:
//     Run test B.
//     Run test C.
//
// Expressed as:
//
//   root := New(A)
//   b := root.Child(B)
//   c := root.Child(C)
//   b.Child(E)
//   b.And(c).Child(D)
//
