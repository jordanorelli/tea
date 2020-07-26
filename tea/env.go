package tea

import (
	"fmt"
	"reflect"
)

type env struct {
	key    string
	value  interface{}
	parent *env
}

func mkenv(test Test) *env {
	var e *env
	return e.save(test)
}

// save looks at the Test t and saves the values of its fields marked with a
// save tag
func (e *env) save(test Test) *env {
	V := reflect.ValueOf(test)
	if V.Type().Kind() == reflect.Ptr {
		V = V.Elem()
	}
	T := V.Type()

	for i := 0; i < T.NumField(); i++ {
		f := T.Field(i)
		if !isSaveField(f) {
			continue
		}

		fv := V.Field(i)
		e = &env{
			key:    f.Name,
			value:  fv.Interface(),
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
			if e.key == f.Name {
				ev := reflect.ValueOf(e.value)
				if ev.Type().AssignableTo(fv.Type()) {
					set = true
					fv.Set(ev)
					break
				}
			}
		}

		if !set {
			return fmt.Errorf("failed to set required field: %q", f.Name)
		}
	}
	return nil
}
