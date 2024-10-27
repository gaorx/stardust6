package sderr

import (
	"maps"
	"reflect"
	"slices"
)

func evalAttrValue(v any) (any, bool) {
	if v == nil {
		return nil, true
	}
	if v1, ok := v.(reflect.Value); ok {
		return evalAttrValue(v1.Interface())
	}
	if v1, ok := v.(func() any); ok {
		if v1 != nil {
			v2 := v1()
			return evalAttrValue(v2)
		}
	}
	vv := reflect.ValueOf(v)
	if vv.IsValid() && vv.Kind() == reflect.Func {
		t := vv.Type()
		if t.NumIn() == 0 && t.NumOut() == 1 {
			retVal := vv.Call([]reflect.Value{})
			return evalAttrValue(retVal[0].Interface())
		} else {
			return nil, false
		}
	}

	return v, true
}

func mergeAttrs(dst map[string]any, e error) {
	unwrappedErrs := UnwrapNested(e)
	slices.Reverse(unwrappedErrs)
	for _, unwrappedErr := range unwrappedErrs {
		if e1, ok := Probe(unwrappedErr); ok {
			maps.Copy(dst, e1.attrs)
		}
	}
}
