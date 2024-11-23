package sdreflect

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"strings"
	"testing"
)

func TestFunc(t *testing.T) {
	is := assert.New(t)

	f1 := func(a int, b string) (int, error) {
		return 0, nil
	}
	f2 := func(a string, b int) int {
		return b * 2
	}
	f3 := func() error {
		return nil
	}
	f4 := func() {}

	is.EqualValues([]reflect.Type{TInt, TString}, Ins(reflect.TypeOf(f1)))
	is.EqualValues([]reflect.Type{TInt, TErr}, Outs(reflect.TypeOf(f1)))
	is.EqualValues([]reflect.Type{TString, TInt}, Ins(reflect.TypeOf(f2)))
	is.EqualValues([]reflect.Type{TInt}, Outs(reflect.TypeOf(f2)))
	res, ok := TrimLastErr(Outs(reflect.TypeOf(f1)))
	is.EqualValues([]reflect.Type{TInt}, res)
	is.True(ok)
	res, ok = TrimLastErr(Outs(reflect.TypeOf(f2)))
	is.EqualValues([]reflect.Type{TInt}, res)
	is.False(ok)
	res, ok = TrimLastErr(Outs(reflect.TypeOf(f3)))
	is.Empty(res)
	is.True(ok)
	res, ok = TrimLastErr(Outs(reflect.TypeOf(f4)))
	is.Empty(res)
	is.False(ok)

	f5 := func(a int, b string) (int, string) {
		return a * 2, strings.Repeat(b, 2)
	}
	inVals := MakeInValues(Ins(reflect.TypeOf(f1)), func(t reflect.Type, i int) reflect.Value {
		if t == TInt {
			return reflect.ValueOf(3)
		}
		if t == TString {
			return reflect.ValueOf("yo")
		}
		return reflect.Value{}
	})
	outVals := reflect.ValueOf(f5).Call(inVals)
	is.EqualValues([]any{6, "yoyo"}, ToInterfaces(outVals...))
}
