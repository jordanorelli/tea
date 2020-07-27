package tea

import (
	"fmt"
	"reflect"
)

type env struct {
	data   map[string]interface{}
	parent *env
}

func mkenv(test Test) *env {
	var e *env
	return e.save(test)
}

// save looks at the Test t and saves the values of its fields marked with a
// save tag. All of the fields for that tests are stored together as a data
// layer.
func (e *env) save(test Test) *env {
	V := reflect.ValueOf(test)
	if V.Type().Kind() == reflect.Ptr {
		V = V.Elem()
	}
	T := V.Type()

	saved := make(map[string]interface{})
	for i := 0; i < T.NumField(); i++ {
		f := T.Field(i)
		if !isSaveField(f) {
			continue
		}

		fv := V.Field(i)
		saved[f.Name] = fv.Interface()
	}
	if len(saved) > 0 {
		return &env{
			data:   saved,
			parent: e,
		}
	}

	return e
}

func (e *env) load(dest Test) error {
	destV := reflect.ValueOf(dest).Elem()
	destT := destV.Type()

	for i := 0; i < destT.NumField(); i++ {
		f := destT.Field(i)
		if !isLoadField(f) {
			continue
		}
		fv := destV.Field(i)

		set := false
		for e := e; e != nil; e = e.parent {
			v, ok := e.data[f.Name]
			if !ok {
				continue
			}

			ev := reflect.ValueOf(v)
			if ev.Type().AssignableTo(fv.Type()) {
				set = true
				fv.Set(ev)
				break
			}
		}

		if !set {
			return fmt.Errorf("failed to set required field: %q", f.Name)
		}
	}
	return nil
}
