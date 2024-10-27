package sderr

import (
	"errors"
	"reflect"
)

type unwrappable interface {
	Unwrap() error
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}

func As[T any](err error) (T, bool) {
	var target T
	if ok := errors.As(err, &target); ok {
		return target, true
	} else {
		var empty T
		return empty, false
	}
}

func Ensure(v any) error {
	if vv, ok := v.(reflect.Value); ok {
		v = vv.Interface()
	}
	switch err := v.(type) {
	case nil:
		return nil
	case error:
		return err
	case string:
		return Newf(err)
	default:
		return Newf("%v", err)
	}
}

func Probe(err error) (*Error, bool) {
	if e, ok := err.(*Error); ok {
		return e, true
	} else {
		return nil, false
	}
}

func ProbeMulti(err error) (*MultiError, bool) {
	if e, ok := err.(*MultiError); ok {
		return e, true
	} else {
		return nil, false
	}
}

func Attrs(err error) map[string]any {
	if err == nil {
		return nil
	}
	merged := map[string]any{}
	mergeAttrs(merged, err)
	return merged
}

func Attr(err error, k string) (any, bool) {
	if err == nil {
		return nil, false
	}
	for {
		if e, ok := Probe(err); ok {
			if v, ok := e.attrs[k]; ok {
				return v, true
			}
		}
		if err = Unwrap(err); err == nil {
			return nil, false
		}
	}
}

func Unwrap(err error) error {
	if err == nil {
		return nil
	}
	if u, ok := err.(unwrappable); ok {
		return u.Unwrap()
	} else {
		return nil
	}
}

func UnwrapNested(err error) []error {
	if err == nil {
		return nil
	}
	var nested []error
	var unwrap func(error)
	unwrap = func(err1 error) {
		if err1 != nil {
			nested = append(nested, err1)
			if u, ok := err1.(unwrappable); ok {
				unwrap(u.Unwrap())
			}
		}
	}
	unwrap(err)
	return nested
}

func Root(err error) error {
	if err == nil {
		return nil
	}
	if err0 := Unwrap(err); err0 == nil {
		return err
	} else {
		return Root(err0)
	}
}

func RootStack(err error) *Stack {
	if err == nil {
		return nil
	}
	wrapsWithRoot := UnwrapNested(err)
	for i := len(wrapsWithRoot) - 1; i >= 0; i-- {
		if e, ok := Probe(wrapsWithRoot[i]); ok {
			return e.stack
		}
	}
	return nil
}
