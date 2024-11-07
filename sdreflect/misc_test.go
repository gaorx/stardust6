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
