package sderr

import (
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestEvalAttrValue(t *testing.T) {
	is := assert.New(t)
	tuple := func(a any, b bool) lo.Tuple2[any, bool] {
		return lo.T2(a, b)
	}
	eval := func(v any) lo.Tuple2[any, bool] {
		return lo.T2(evalAttrValue(v))
	}

	// ok
	is.Equal(tuple(nil, true), eval(nil))
	is.Equal(tuple(33, true), eval(reflect.ValueOf(33)))
	is.Equal(tuple(44, true), eval(func() any {
		return 44
	}))
	is.Equal(tuple(55, true), eval(func() int {
		return 55
	}))
	is.Equal(tuple(66, true), eval(func() any {
		return reflect.ValueOf(66)
	}))
	is.Equal(tuple(77, true), eval(func() reflect.Value {
		return reflect.ValueOf(77)
	}))
	is.Equal(tuple("ABC", true), eval("ABC"))
	is.Equal(tuple(true, true), eval(true))
	is.Equal(tuple(int64(33), true), eval(int64(33)))
	is.Equal(tuple(33.3, true), eval(33.3))

	// no
	is.Equal(tuple(nil, false), eval(func(int) int {
		return 33
	}))
}
