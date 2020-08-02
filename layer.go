package tea

import "reflect"

type sval struct {
	name string
	val  interface{}
}

// layerData is an ordered list of key-value pairs. We preserve the ordering of
// the fields from the structs that produced the layer so that viewing a list
// of layers shows each layer from the sam struct in the same value-order,
// which is the order those fields appear in their originating stuct (and not,
// like, alphabetical order). Although a bit more tedious to work with
// internally, this is intended to make it possible to write more descriptive,
// easily understood PlanError messages.
type layerData []sval

func (l layerData) get(key string) (val interface{}, present bool) {
	for _, pair := range l {
		if pair.name == key {
			return pair.val, true
		}
	}
	return nil, false
}

type layer struct {
	origin *xnode
	saved  layerData
}

func makeLayerData(test Test) (layerData, error) {
	V := reflect.ValueOf(test)
	if V.Type().Kind() == reflect.Ptr {
		V = V.Elem()
	}
	T := V.Type()

	if T.Kind() != reflect.Struct {
		return nil, nil
	}

	fields, err := getSaveFields(T)
	if err != nil {
		return nil, err
	}
	if len(fields) == 0 {
		// is this weird? maybe this is weird.
		return nil, nil
	}

	data := make(layerData, 0, len(fields))
	for _, f := range fields {
		fv := V.FieldByName(f.Name).Interface()
		data = append(data, sval{name: f.Name, val: fv})
	}
	return data, nil
}
