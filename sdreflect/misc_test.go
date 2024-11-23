package sdreflect

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestRootValueOf(t *testing.T) {
	is := assert.New(t)
	caseVal0 := reflect.ValueOf(reflect.ValueOf(reflect.ValueOf(3)))
	caseVals := []any{
		3,
		reflect.ValueOf(3),
		reflect.ValueOf(reflect.ValueOf(3)),
		reflect.ValueOf(reflect.ValueOf(reflect.ValueOf(3))),
		&caseVal0,
	}
	for _, v0 := range caseVals {
		is.Equal(3, RootValueOf(v0).Interface())
	}
	is.Panics(func() {
		RootValueOf((*reflect.Value)(nil))
	})
}

func TestDeref(t *testing.T) {
	is := assert.New(t)

	var a = 3
	var pa = &a
	var ppa = &pa
	is.Equal(reflect.ValueOf(a).Interface(), Deref(reflect.ValueOf(a)).Interface())
	is.Equal(reflect.ValueOf(a).Interface(), Deref(reflect.ValueOf(pa)).Interface())
	is.Equal(reflect.ValueOf(a).Interface(), Deref(reflect.ValueOf(ppa)).Interface())
}

func TestTypesOf(t *testing.T) {
	is := assert.New(t)
	is.EqualValues(
		[]reflect.Type{TString, TInt, TBool},
		ToTypes("a", 1, true),
	)
	is.EqualValues(
		[]reflect.Type{TString, TInt, TBool},
		TypesOf(ToValues("a", 1, true)),
	)
}

func TestToInterfaces(t *testing.T) {
	is := assert.New(t)
	is.EqualValues(
		[]any{"a", 1, true},
		ToInterfaces(ToValues("a", 1, true)...),
	)
}
