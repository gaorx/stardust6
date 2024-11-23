package sdreflect

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestIs(t *testing.T) {
	is := assert.New(t)

	is.True(IsStruct(reflect.TypeOf(struct{}{})))
	is.True(IsStructPtr(reflect.TypeOf(&struct{}{})))

	is.True(IsSlice(reflect.TypeOf([]int{}), nil))
	is.True(IsSlice(reflect.TypeOf([]int{}), reflect.TypeOf(0)))
	is.False(IsSlice(reflect.TypeOf([]int{}), reflect.TypeOf("")))

	is.True(IsSliceLike(reflect.TypeOf([]int{}), nil))
	is.True(IsSliceLike(reflect.TypeOf([]int{}), reflect.TypeOf(0)))
	is.False(IsSliceLike(reflect.TypeOf([]int{}), reflect.TypeOf("")))
	is.True(IsSliceLike(reflect.TypeOf([3]int{}), nil))
	is.True(IsSliceLike(reflect.TypeOf([3]int{}), reflect.TypeOf(0)))
	is.False(IsSliceLike(reflect.TypeOf([3]int{}), reflect.TypeOf("")))

	is.True(IsMap(reflect.TypeOf(map[int]string{}), nil, nil))
	is.True(IsMap(reflect.TypeOf(map[int]string{}), reflect.TypeOf(0), nil))
	is.True(IsMap(reflect.TypeOf(map[int]string{}), nil, reflect.TypeOf("")))
	is.True(IsMap(reflect.TypeOf(map[int]string{}), reflect.TypeOf(0), reflect.TypeOf("")))
	is.False(IsMap(reflect.TypeOf(map[int]string{}), nil, reflect.TypeOf(0)))

	is.True(IsTypes1[int](ToTypes(1)))
	is.True(IsTypes2[int, string](ToTypes(1, "a")))
	is.True(IsTypes3[int, string, bool](ToTypes(1, "a", true)))
	is.False(IsTypes4[int, string, bool, error](ToTypes(1, "a", true, errors.New("xx"))))

	is.True(IsAssignableTypes1[int](ToTypes(1)))
	is.True(IsAssignableTypes2[int, string](ToTypes(1, "a")))
	is.True(IsAssignableTypes3[int, string, bool](ToTypes(1, "a", true)))
	is.True(IsAssignableTypes4[int, string, bool, error](ToTypes(1, "a", true, errors.New("xx"))))
}
