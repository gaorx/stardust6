package sderr

import (
	"errors"
	"reflect"
)

type unwrappable interface {
	Unwrap() error
}

// Is 相当于errors.Is
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// As errors.As的范型版
func As[T any](err error) (T, bool) {
	var target T
	if ok := errors.As(err, &target); ok {
		return target, true
	} else {
		var empty T
		return empty, false
	}
}

// Ensure 将任何值转换成一个error
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

// Probe 探测一个error是否是*sderr.Error，并且返回*sderr.Error
func Probe(err error) (*Error, bool) {
	if e, ok := err.(*Error); ok {
		return e, e != nil
	} else {
		return nil, false
	}
}

// ProbeMulti 探测一个error是否是*sderr.MultiError，并且返回*sderr.MultiError
func ProbeMulti(err error) (*MultiError, bool) {
	if e, ok := err.(*MultiError); ok {
		return e, e != nil
	} else {
		return nil, false
	}
}

// Attrs 返回一个error的所有attrs，包括嵌套的多层wrap中的attrs
func Attrs(err error) map[string]any {
	if err == nil {
		return nil
	}
	merged := map[string]any{}
	mergeAttrs(merged, err)
	return merged
}

// Attr 返回一个error的某个attr，包括嵌套的多层wrap中的attr
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

// Unwrap 返回一个error的下一层error
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

// UnwrapNested 返回一个error的所有嵌套的下一层error
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

// Root 返回一个error的根cause
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

// RootStack 返回一个error的根cause的stacktrace
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
