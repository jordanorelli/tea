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

	e = e.match(dest)
	if e == nil {
		return fmt.Errorf("failed to find a matching environment")
	}

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

func (e *env) match(dest Test) *env {
	destV := reflect.ValueOf(dest).Elem()
	destT := destV.Type()

	required := getMatchFields(destT)
	if len(required) == 0 {
		return e
	}

	var last *env
	var leaf *env

	for e := e; e != nil; e = e.parent {
		present := make([]reflect.StructField, 0, len(required))

		for _, f := range required {
			ev, ok := e.data[f.Name]
			if !ok {
				break
			}
			if reflect.TypeOf(ev).AssignableTo(f.Type) {
				present = append(present, f)
			}
		}

		// all required fields are present in this layer
		if len(present) == len(required) {
			// check that the values in the env match the values that were
			// asked for.
			matched := make(map[string]interface{})
			for _, f := range required {
				fv := destV.FieldByName(f.Name)
				if fv.Interface() == e.data[f.Name] {
					matched[f.Name] = e.data[f.Name]
				}
			}

			// all required match conditions are met
			if len(matched) == len(required) {
				if leaf == nil {
					// if this is the first matched layer, it is the leaf of the
					// resultant env.
					leaf = e
					last = leaf
				} else {
					// otherwise we keep this layer, since it matched our match
					// requirements. Another layer already did, but there may
					// be other things in the layer we want to keep.
					last.parent = e
					last = e
				}
			}
		} else {
			// the required fields do not exist in the layer, so this layer
			// does not conflict with the match requirement.
			if leaf != nil {
				// since we have a leaf node, we have found a matching layer,
				// and since this layer does not conflict, we keep it.
				last.parent = e
				last = e
			}
		}
	}

	return leaf
}
