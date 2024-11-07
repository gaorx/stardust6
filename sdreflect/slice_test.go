package sdreflect

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestIsSliceLikeValue(t *testing.T) {
	is := assert.New(t)
	a := [3]int{1, 2, 3}
	is.True(IsSliceLikeValue(reflect.ValueOf(a)))
	is.True(IsSliceLikeValue(reflect.ValueOf(a[:])))
}

func TestIsSliceLikeType(t *testing.T) {
	is := assert.New(t)
	a := [3]int{1, 2, 3}
	is.True(IsSliceLikeType(reflect.TypeOf(a)))
	is.True(IsSliceLikeType(reflect.TypeOf(a[:])))
}

func TestForEach(t *testing.T) {
	is := assert.New(t)
	a := [3]int{3, 2, 1}
	ForEach(reflect.ValueOf(a), func(elem reflect.Value, i, n int) {
		is.True(a[i] == elem.Interface().(int))
		is.Equal(len(a), n)
	})
	ForEach(reflect.ValueOf(a[:]), func(elem reflect.Value, i, n int) {
		is.True(a[i] == elem.Interface().(int))
		is.Equal(len(a), n)
	})
}
