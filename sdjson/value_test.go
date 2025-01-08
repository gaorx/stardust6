package sdjson

import (
	"encoding/json"
	"math"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestV(t *testing.T) {
	assert.Equal(t, "hello", V("hello").Interface())
	assert.Equal(t, nil, V(nil).Interface())
	assert.True(t, V(nil).IsNil())
}

// Get

func TestValueField(t *testing.T) {
	v, err := UnmarshalValueString(`{
		"k1": {
			"k2": 33,
			"k3": "mm",
			"k4": true
		}
	}`)
	assert.NoError(t, err)
	assert.Equal(t, 33, v.Get("k1", "k2").AsIntOr(0))
	assert.Equal(t, "mm", v.Get("k1").Get("k3").AsStringOr(""))
	assert.Equal(t, "true", v.Get("k1", "k4").AsStringOr(""))
	assert.Equal(t, "not_found", v.Get("k1").Get("k5").AsStringOr("not_found"))
	assert.Equal(t, "not_found", v.Get("k2").AsStringOr("not_found"))
}

func TestValue_At(t *testing.T) {
	v, err := UnmarshalValueString(`["a", 3, {"k1":"v1"}]`)
	assert.NoError(t, err)
	assert.Equal(t, "a", v.At(0).AsStringOr(""))
	assert.Equal(t, "3", v.At(1).AsStringOr(""))
	assert.Equal(t, "v1", v.At(2).Get("k1").AsStringOr(""))
}

// ToXXX

func TestValueToBool(t *testing.T) {
	// bool
	newfr(V(true).ToBool(false)).with(t).noErr().equal(true)
	newfr(V(false).ToBool(false)).with(t).noErr().equal(false)

	// other
	newfr(V(0).ToBool(false)).with(t).hasErr()
	newfr(V("true").ToBool(false)).with(t).hasErr()
	newfr(V(1.3).ToBool(false)).with(t).hasErr()
}

func TestValueToString(t *testing.T) {
	// string
	newfr(V("xx").ToString(false)).with(t).noErr().equal("xx")

	// other
	newfr(V(0).ToString(false)).with(t).hasErr()
	newfr(V(true).ToString(false)).with(t).hasErr()
	newfr(V(1.0).ToString(false)).with(t).hasErr()
}

func TestValueToInt(t *testing.T) {
	// other
	newfr(V(nil).ToInt(false)).with(t).hasErr()
	newfr(V(true).ToInt(false)).with(t).hasErr()
	newfr(V("0").ToInt(false)).with(t).hasErr()
	newfr(V(Object{}).ToInt(false)).with(t).hasErr()
	newfr(V(Array{}).ToInt(false)).with(t).hasErr()

	// number
	newfr(V(3).ToInt(false)).with(t).noErr().equal(3)
	newfr(V(int8(3)).ToInt(false)).with(t).noErr().equal(3)
	newfr(V(int16(3)).ToInt(false)).with(t).noErr().equal(3)
	newfr(V(int32(3)).ToInt(false)).with(t).noErr().equal(3)
	newfr(V(int64(3)).ToInt(false)).with(t).noErr().equal(3)
	newfr(V(uint(3)).ToInt(false)).with(t).noErr().equal(3)
	newfr(V(uint8(3)).ToInt(false)).with(t).noErr().equal(3)
	newfr(V(uint16(3)).ToInt(false)).with(t).noErr().equal(3)
	newfr(V(uint32(3)).ToInt(false)).with(t).noErr().equal(3)
	newfr(V(uint64(3)).ToInt(false)).with(t).noErr().equal(3)
	newfr(V(3.3).ToInt(false)).with(t).noErr().equal(3)
	newfr(V(float32(3.3)).ToInt(false)).with(t).noErr().equal(3)
	newfr(V(json.Number("3.3")).ToInt(false)).with(t).noErr().equal(3)
	newfr(V(json.Number("3")).ToInt(false)).with(t).noErr().equal(3)
}

func TestValueToUint(t *testing.T) {
	// other
	newfr(V(nil).ToUint(false)).with(t).hasErr()
	newfr(V(true).ToUint(false)).with(t).hasErr()
	newfr(V("0").ToUint(false)).with(t).hasErr()
	newfr(V(Object{}).ToUint(false)).with(t).hasErr()
	newfr(V(Array{}).ToUint(false)).with(t).hasErr()

	// number
	newfr(V(3).ToUint(false)).with(t).noErr().equal(uint(3))
	newfr(V(int8(3)).ToUint(false)).with(t).noErr().equal(uint(3))
	newfr(V(int16(3)).ToUint(false)).with(t).noErr().equal(uint(3))
	newfr(V(int32(3)).ToUint(false)).with(t).noErr().equal(uint(3))
	newfr(V(int64(3)).ToUint(false)).with(t).noErr().equal(uint(3))
	newfr(V(uint(3)).ToUint(false)).with(t).noErr().equal(uint(3))
	newfr(V(uint8(3)).ToUint(false)).with(t).noErr().equal(uint(3))
	newfr(V(uint16(3)).ToUint(false)).with(t).noErr().equal(uint(3))
	newfr(V(uint32(3)).ToUint(false)).with(t).noErr().equal(uint(3))
	newfr(V(uint64(3)).ToUint(false)).with(t).noErr().equal(uint(3))
	newfr(V(3.3).ToUint(false)).with(t).noErr().equal(uint(3))
	newfr(V(float32(3.3)).ToUint(false)).with(t).noErr().equal(uint(3))
	newfr(V(json.Number("3.3")).ToUint(false)).with(t).noErr().equal(uint(3))
	newfr(V(json.Number("3")).ToUint(false)).with(t).noErr().equal(uint(3))
}

func TestValueToFloat64(t *testing.T) {
	// other
	newfr(V(nil).ToUint(false)).with(t).hasErr()
	newfr(V(true).ToUint(false)).with(t).hasErr()
	newfr(V("0").ToUint(false)).with(t).hasErr()
	newfr(V(Object{}).ToUint(false)).with(t).hasErr()
	newfr(V(Array{}).ToUint(false)).with(t).hasErr()

	// number
	newfr(V(3).ToFloat64(false)).with(t).noErr().equal(3.0)
	newfr(V(int8(3)).ToFloat64(false)).with(t).noErr().equal(3.0)
	newfr(V(int16(3)).ToFloat64(false)).with(t).noErr().equal(3.0)
	newfr(V(int32(3)).ToFloat64(false)).with(t).noErr().equal(3.0)
	newfr(V(int64(3)).ToFloat64(false)).with(t).noErr().equal(3.0)
	newfr(V(uint(3)).ToFloat64(false)).with(t).noErr().equal(3.0)
	newfr(V(uint8(3)).ToFloat64(false)).with(t).noErr().equal(3.0)
	newfr(V(uint16(3)).ToFloat64(false)).with(t).noErr().equal(3.0)
	newfr(V(uint32(3)).ToFloat64(false)).with(t).noErr().equal(3.0)
	newfr(V(uint64(3)).ToFloat64(false)).with(t).noErr().equal(3.0)
	newfr(V(3.3).ToFloat64(false)).with(t).noErr().equal(3.3)
	newfr(V(float32(3.3)).ToFloat64(false)).with(t).noErr().equalFloat64(3.3)
	newfr(V(json.Number("3.3")).ToFloat64(false)).with(t).noErr().equal(3.3)
	newfr(V(json.Number("3")).ToFloat64(false)).with(t).noErr().equal(3.0)
}

func TestValueToObject(t *testing.T) {
	// other
	newfr(V(nil).ToObject(false)).with(t).hasErr()
	newfr(V(true).ToObject(false)).with(t).hasErr()
	newfr(V("0").ToObject(false)).with(t).hasErr()
	newfr(V(0.1).ToObject(false)).with(t).hasErr()
	newfr(V(Array{}).ToObject(false)).with(t).hasErr()

	// Objects
	newfr(V(Object{}).ToObject(false)).with(t).isObject()
	newfr(V(Object(nil)).ToObject(false)).with(t).isNil()
	newfr(V(map[string]any{"k1": "v1"}).ToObject(false)).with(t).isObject().deepEqual(Object{"k1": "v1"})
	newfr(V(map[string]any(nil)).ToObject(false)).with(t).isNil()
	newfr(V(map[string]int{"k1": 0}).ToObject(false)).with(t).isObject().deepEqual(Object{"k1": 0})
	newfr(V(map[string]bool(nil)).ToObject(false)).with(t).isNil()
}

func TestValueToArray(t *testing.T) {
	//// other
	newfr(V(nil).ToArray(false)).with(t).hasErr()
	newfr(V(true).ToArray(false)).with(t).hasErr()
	newfr(V("0").ToArray(false)).with(t).hasErr()
	newfr(V(0.1).ToArray(false)).with(t).hasErr()
	newfr(V(Object{}).ToArray(false)).with(t).hasErr()
	//
	//// Array
	newfr(V(Array{}).ToArray(false)).with(t).isArray()
	newfr(V(Array{"a", 1}).ToArray(false)).with(t).isArray().deepEqual(Array{"a", 1})
	newfr(V(Array(nil)).ToArray(false)).with(t).isNil()
	newfr(V([]any{"a"}).ToArray(false)).with(t).isArray().deepEqual(Array{"a"})
	newfr(V([]any(nil)).ToArray(false)).with(t).isNil()
	newfr(V([]string{"a"}).ToArray(false)).with(t).isArray().deepEqual(Array{"a"})
	newfr(V([3]int{33}).ToArray(false)).with(t).isArray().deepEqual(Array{33, 0, 0})
}

func TestValueToAny(t *testing.T) {
	type person struct {
		Name string `json:"name"`
	}
	var p1 person
	assert.True(t, V(Object{"name": "xx"}).To(&p1, false))
	assert.Equal(t, "xx", p1.Name)

	var p2 person
	assert.True(t, V(person{"yy"}).To(&p2, false))
	assert.Equal(t, "yy", p2.Name)

	var p3 person
	assert.True(t, V(person{"zz"}).To(&p3, false))
	assert.Equal(t, "zz", p3.Name)

	var p4 Object
	assert.True(t, V(person{Name: "oo"}).To(&p4, false))
	assert.True(t, reflect.DeepEqual(p4, Object{"name": "oo"}))
}

// TryXXX

func TestValueToBoolAs(t *testing.T) {
	// bool
	newfr(V(true).ToBool(true)).with(t).noErr().equal(true)
	newfr(V(false).ToBool(true)).with(t).noErr().equal(false)

	// int
	newfr(V(0).ToBool(true)).with(t).noErr().equal(false)
	newfr(V(1).ToBool(true)).with(t).noErr().equal(true)
	newfr(V(2).ToBool(true)).with(t).noErr().equal(true)

	// uint
	// int
	newfr(V(uint(0)).ToBool(true)).with(t).noErr().equal(false)
	newfr(V(uint(1)).ToBool(true)).with(t).noErr().equal(true)
	newfr(V(uint(2)).ToBool(true)).with(t).noErr().equal(true)

	// string
	newfr(V("true").ToBool(true)).with(t).noErr().equal(true)
	newfr(V("false").ToBool(true)).with(t).noErr().equal(false)

	// float
	newfr(V(0.0).ToBool(true)).with(t).noErr().equal(false)
	newfr(V(1.0).ToBool(true)).with(t).noErr().equal(true)
	newfr(V(3.3).ToBool(true)).with(t).noErr().equal(true)

	// other
	newfr(V(Object{}).ToBool(true)).with(t).hasErr()
	newfr(V(Array{}).ToBool(true)).with(t).hasErr()
}

func TestValueToStringAs(t *testing.T) {
	// bool
	newfr(V(true).ToString(true)).with(t).noErr().equal("true")
	newfr(V(false).ToString(true)).with(t).noErr().equal("false")

	// string
	newfr(V("xx").ToString(true)).with(t).noErr().equal("xx")

	// int
	newfr(V(-33).ToString(true)).with(t).noErr().equal("-33")

	// uint
	newfr(V(uint(33)).ToString(true)).with(t).noErr().equal("33")

	// float64
	newfr(V(3.3).ToString(true)).with(t).noErr().equal("3.3")

	// object
	newfr(V(Object{}).ToString(true)).with(t).hasErr()

	// array
	newfr(V(Array{}).ToString(true)).with(t).hasErr()
}

func TestValueToIntAs(t *testing.T) {
	// bool
	newfr(V(true).ToInt(true)).with(t).noErr().equal(1)
	newfr(V(false).ToInt(true)).with(t).noErr().equal(0)

	// string
	newfr(V("33").ToInt(true)).with(t).noErr().equal(33)

	// int
	newfr(V(-33).ToInt(true)).with(t).noErr().equal(-33)

	// uint
	newfr(V(uint(33)).ToInt(true)).with(t).noErr().equal(33)

	// float64
	newfr(V(3.3).ToInt(true)).with(t).noErr().equal(3)

	// object
	newfr(V(Object{}).ToInt(true)).with(t).hasErr()

	// array
	newfr(V(Array{}).ToInt(true)).with(t).hasErr()
}

func TestValueToUintAs(t *testing.T) {
	// bool
	newfr(V(true).ToUint(true)).with(t).noErr().equal(uint(1))
	newfr(V(false).ToUint(true)).with(t).noErr().equal(uint(0))

	// string
	newfr(V("33").ToUint(true)).with(t).noErr().equal(uint(33))

	// int
	x := -33
	newfr(V(-33).ToUint(true)).with(t).noErr().equal(uint(x))

	// uint
	newfr(V(uint(33)).ToUint(true)).with(t).noErr().equal(uint(33))

	// float64
	newfr(V(3.3).ToUint(true)).with(t).noErr().equal(uint(3))

	// object
	newfr(V(Object{}).ToUint(true)).with(t).hasErr()

	// array
	newfr(V(Array{}).ToUint(true)).with(t).hasErr()
}

func TestValueToFloat64As(t *testing.T) {
	// bool
	newfr(V(true).ToFloat64(true)).with(t).noErr().equal(1.0)
	newfr(V(false).ToFloat64(true)).with(t).noErr().equal(0.0)

	// string
	newfr(V("33.3").ToFloat64(true)).with(t).noErr().equal(33.3)

	// int
	newfr(V(-33).ToFloat64(true)).with(t).noErr().equal(-33.0)

	// uint
	newfr(V(uint(33)).ToFloat64(true)).with(t).noErr().equal(33.0)

	// float64
	newfr(V(3.3).ToFloat64(true)).with(t).noErr().equal(3.3)

	// object
	newfr(V(Object{}).ToFloat64(true)).with(t).hasErr()

	// array
	newfr(V(Array{}).ToFloat64(true)).with(t).hasErr()
}

func TestValueToObjectAs(t *testing.T) {
	type person struct {
		Name string `json:"name"`
	}
	newfr(V(&person{Name: "xx"}).ToObject(true)).with(t).noErr().deepEqual(Object{"name": "xx"})
	newfr(V(person{Name: "yy"}).ToObject(true)).with(t).noErr().deepEqual(Object{"name": "yy"})
}

type funcReturn struct {
	v  any
	ok bool
	t  *testing.T
}

func newfr(v any, ok bool) *funcReturn {
	return &funcReturn{v, ok, nil}
}

func (fr *funcReturn) with(t *testing.T) *funcReturn {
	fr.t = t
	return fr
}

func (fr *funcReturn) hasErr() *funcReturn {
	assert.False(fr.t, fr.ok)
	return fr
}

func (fr *funcReturn) noErr() *funcReturn {
	assert.True(fr.t, fr.ok)
	return fr
}

func (fr *funcReturn) equal(expected any) *funcReturn {
	assert.Equal(fr.t, expected, fr.v)
	return fr
}

func (fr *funcReturn) equalFloat64(expected float64) *funcReturn {
	assert.IsType(fr.t, 0.0, fr.v)
	v1 := fr.v.(float64)
	assert.True(fr.t, math.Abs(expected-v1) < 0.0000001)
	return fr
}

func (fr *funcReturn) isObject() *funcReturn {
	assert.IsType(fr.t, Object(nil), fr.v)
	return fr
}

func (fr *funcReturn) isArray() *funcReturn {
	assert.IsType(fr.t, Array(nil), fr.v)
	return fr
}

func (fr *funcReturn) isNil() *funcReturn {
	assert.Nil(fr.t, fr.v)
	return fr
}

func (fr *funcReturn) deepEqual(expected any) *funcReturn {
	assert.True(fr.t, reflect.DeepEqual(expected, fr.v))
	return fr
}
