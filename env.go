package tea

import (
	"bytes"
	"fmt"
	"reflect"
	"sort"
	"strings"
)

type env struct {
	data   map[string]interface{}
	parent *env
}

func (e *env) String() string {
	var buf bytes.Buffer
	e.pretty(&buf)
	return buf.String()
}

func (e *env) pretty(buf *bytes.Buffer) {
	if e.parent != nil {
		e.parent.pretty(buf)
		buf.WriteRune(' ')
	}

	var parts []string
	for k, v := range e.data {
		if s, ok := v.(string); ok {
			parts = append(parts, fmt.Sprintf("[%s=%q]", k, s))
		} else {
			parts = append(parts, fmt.Sprintf("[%s=%v]", k, v))
		}
	}
	sort.Strings(parts)

	fmt.Fprintf(buf, "{%s}", strings.Join(parts, ", "))
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

	if T.Kind() != reflect.Struct {
		return e
	}

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

	e, err := e.match(dest)
	if err != nil {
		return fmt.Errorf("match failed: %w", err)
	}

	for i := 0; i < destT.NumField(); i++ {
		f := destT.Field(i)
		if !isLoadField(f) {
			continue
		}
		fv := destV.Field(i)
		if !fv.IsZero() {
			// the value is already populated, so we don't want to overwrite
			// it.
			continue
		}

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
			return fmt.Errorf("%w: failed to set required field: %q", PlanError, f.Name)
		}
	}
	return nil
}

func (e *env) match(dest Test) (*env, error) {
	destV := reflect.ValueOf(dest).Elem()
	destT := destV.Type()

	required := getMatchFields(destT)
	if len(required) == 0 {
		return e, nil
	}

	var (
		last                  *env
		leaf                  *env
		foundWithWrongType    = make(map[string]bool)
		foundWithMatchingType = make(map[string]bool)
		foundWithWrongValue   = make(map[string]bool)
		foundWithCorrectValue = make(map[string]bool)
	)

	keep := func(v *env) {
		if leaf == nil {
			leaf = &env{data: v.data}
			last = leaf
		} else {
			next := &env{data: v.data, parent: last}
			last = next
		}
	}

	for e := e; e != nil; e = e.parent {
		present := make([]reflect.StructField, 0, len(required))

		for _, f := range required {
			ev, ok := e.data[f.Name]
			if !ok {
				break
			}
			if reflect.TypeOf(ev).AssignableTo(f.Type) {
				foundWithMatchingType[f.Name] = true
				present = append(present, f)
			} else {
				foundWithWrongType[f.Name] = true
			}
		}

		// all required fields are present in this layer
		if len(present) == len(required) {
			// check that the values in the env match the values that were
			// asked for.
			matched := make(map[string]interface{})
			wrongVal := make(map[string]bool)

			for _, f := range required {
				fv := destV.FieldByName(f.Name)
				if fv.Interface() == e.data[f.Name] {
					foundWithCorrectValue[f.Name] = true
					matched[f.Name] = e.data[f.Name]
				} else {
					foundWithWrongValue[f.Name] = true
					wrongVal[f.Name] = true
				}
			}

			if len(wrongVal) > 0 {
				continue
			}

			// all required match conditions are met
			if len(matched) == len(required) {
				keep(e)
			}
		} else {
			// the required fields do not exist in the layer, so this layer
			// does not conflict with the match requirement.
			if leaf != nil {
				// since we have a leaf node, we have already found a matching
				// layer, and since this layer does not conflict, we keep it.
				keep(e)
			}
		}
	}

	if leaf == nil {
		var notFound []string
		for _, f := range required {
			if !foundWithMatchingType[f.Name] && !foundWithWrongType[f.Name] {
				notFound = append(notFound, f.Name)
			}
		}
		switch len(notFound) {
		case 0:
			break
		case 1:
			return nil, fmt.Errorf("%w: missing required field: %q", PlanError, notFound[0])
		default:
			return nil, fmt.Errorf("%w: missing %d required fields: %s", PlanError, len(notFound), notFound)
		}
		for f, _ := range foundWithWrongType {
			if !foundWithMatchingType[f] {
				return nil, fmt.Errorf("%w: field %s was only found with unmatching types", PlanError, f)
			}
		}
		for f, _ := range foundWithWrongValue {
			if !foundWithCorrectValue[f] {
				return nil, fmt.Errorf("%w: field %s was only found with unmatching values", RunError, f)
			}
		}
		return nil, fmt.Errorf("%w: required match fields not encountered on the same layer", PlanError)
	}

	return leaf, nil
}
