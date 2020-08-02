package tea

import (
	"math/rand"
	"testing"
)

func TestLayerData(t *testing.T) {
	t.Run("non-struct tests produce empty layers", func(t *testing.T) {
		lr, err := makeLayerData(Pass)
		if len(lr) != 0 {
			t.Errorf("expected a nil layer but saw %v instead", lr)
		}
		if err != nil {
			t.Errorf("expected no error from lay but saw %v instead", err)
		}
	})

	t.Run("save tags on unexported fields are plan errors", func(t *testing.T) {
		type T struct {
			Passing
			count int `tea:"save"`
		}
		_, err := makeLayerData(T{})
		assertError(t, true, err, PlanError)
	})

	t.Run("mixed exported/unexported fields still an error", func(t *testing.T) {
		type T struct {
			Passing
			count int    `tea:"save"`
			Bar   string `tea:"save"`
		}
		_, err := makeLayerData(T{})
		assertError(t, true, err, PlanError)
	})

	t.Run("mixed exported/unexported fields still an error", func(t *testing.T) {
		type T struct {
			Passing
			count int    `tea:"save"`
			bar   string `tea:"save"`
		}
		_, err := makeLayerData(T{})
		assertError(t, true, err, PlanError)
	})

	t.Run("save one int", func(t *testing.T) {
		type T struct {
			Passing
			Count int `tea:"save"`
		}
		test := T{Count: rand.Int()}
		data, err := makeLayerData(test)
		assertNoError(t, true, err)
		if len(data) == 0 {
			t.Fatalf("expected nonempty layer, saw empty layer instead")
		}
		if v, ok := data.get("Count"); !ok {
			t.Errorf("layer data is missing expected field Count")
		} else {
			if v != test.Count {
				t.Errorf("layer data expected Count value of %d but saw %d instead", test.Count, v)
			}
		}
	})

	t.Run("an int and a string", func(t *testing.T) {
		type T struct {
			Passing
			Count int    `tea:"save"`
			Name  string `tea:"save"`
		}
		test := T{Count: rand.Int(), Name: rstring(8)}
		data, err := makeLayerData(test)
		assertNoError(t, true, err)
		if len(data) == 0 {
			t.Fatalf("expected nonempty layer, saw empty layer instead")
		}
		if v, ok := data.get("Count"); !ok {
			t.Errorf("layer data is missing expected field Count")
		} else {
			if v != test.Count {
				t.Errorf("layer data expected Count value of %d but saw %d instead", test.Count, v)
			}
		}
		if v, ok := data.get("Name"); !ok {
			t.Errorf("layer data is missing expected field Count")
		} else {
			if v != test.Name {
				t.Errorf("layer data expected Name value of %s but saw %s instead", test.Name, v)
			}
		}
	})
}
