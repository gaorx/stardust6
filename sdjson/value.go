package sdjson

import (
	"encoding/json"
	"github.com/gaorx/stardust6/sderr"
)

// Value 描述JSON中的一个值，它可能是object, array, string, bool, number, null
type Value struct {
	v any
}

// V 创建一个Value
func V(v any) Value {
	return Value{unbox(v)}
}

// MarshalJSON 将Value转换为JSON，实现json.Marshaler接口
func (v Value) MarshalJSON() ([]byte, error) {
	raw, err := json.Marshal(v.v)
	if err != nil {
		return nil, sderr.Wrapf(err, "sdjson marshal value to json")
	}
	return raw, nil
}

// UnmarshalJSON 将JSON转换为Value，实现json.Unmarshaler接口
func (v *Value) UnmarshalJSON(raw []byte) error {
	var v0 any
	err := json.Unmarshal(raw, &v0)
	if err != nil {
		return sderr.Wrapf(err, "sdjson unmarshal json to value")
	}
	v.v = v0
	return nil
}

// Interface 返回Value的原始值
func (v Value) Interface() any {
	return v.v
}

// IsNil 判断Value是否为nil
func (v Value) IsNil() bool {
	return v.v == nil
}

// Get 如果值是object时获取Value中的一个属性，如果object是嵌套的，可以使用subKeys返回内层的值
func (v Value) Get(k string, subKeys ...string) Value {
	if len(subKeys) <= 0 {
		return v.get(k)
	}
	r := v.get(k)
	for _, subKey := range subKeys {
		r = r.get(subKey)
	}
	return r
}

// At 如果值是array时获取Value中的一个元素
func (v Value) At(i int) Value {
	return v.AsArrayOr(nil).At(i)
}

func (v Value) ToBool(as bool) (bool, bool)       { return ToBool(v, as) }
func (v Value) ToString(as bool) (string, bool)   { return ToString(v, as) }
func (v Value) ToInt(as bool) (int, bool)         { return ToInt[int](v, as) }
func (v Value) ToInt8(as bool) (int8, bool)       { return ToInt[int8](v, as) }
func (v Value) ToInt16(as bool) (int16, bool)     { return ToInt[int16](v, as) }
func (v Value) ToInt32(as bool) (int32, bool)     { return ToInt[int32](v, as) }
func (v Value) ToInt64(as bool) (int64, bool)     { return ToInt[int64](v, as) }
func (v Value) ToUint(as bool) (uint, bool)       { return ToUint[uint](v, as) }
func (v Value) ToUint8(as bool) (uint8, bool)     { return ToUint[uint8](v, as) }
func (v Value) ToUint16(as bool) (uint16, bool)   { return ToUint[uint16](v, as) }
func (v Value) ToUint32(as bool) (uint32, bool)   { return ToUint[uint32](v, as) }
func (v Value) ToUint64(as bool) (uint64, bool)   { return ToUint[uint64](v, as) }
func (v Value) ToFloat32(as bool) (float32, bool) { return ToFloat[float32](v, as) }
func (v Value) ToFloat64(as bool) (float64, bool) { return ToFloat[float64](v, as) }
func (v Value) ToObject(as bool) (Object, bool)   { return ToObject(v, as) }
func (v Value) ToArray(as bool) (Array, bool)     { return ToArray(v, as) }
func (v Value) To(ptr any, as bool) bool          { return converter.ToAny(v, as, ptr) }

func (v Value) AsBool() bool       { return AsBool(v) }
func (v Value) AsString() string   { return AsString(v) }
func (v Value) AsInt() int         { return AsInt[int](v) }
func (v Value) AsInt8() int8       { return AsInt[int8](v) }
func (v Value) AsInt16() int16     { return AsInt[int16](v) }
func (v Value) AsInt32() int32     { return AsInt[int32](v) }
func (v Value) AsInt64() int64     { return AsInt[int64](v) }
func (v Value) AsUint() uint       { return AsUint[uint](v) }
func (v Value) AsUint8() uint8     { return AsUint[uint8](v) }
func (v Value) AsUint16() uint16   { return AsUint[uint16](v) }
func (v Value) AsUint32() uint32   { return AsUint[uint32](v) }
func (v Value) AsUint64() uint64   { return AsUint[uint64](v) }
func (v Value) AsFloat32() float32 { return AsFloat[float32](v) }
func (v Value) AsFloat64() float64 { return AsFloat[float64](v) }
func (v Value) AsObject() Object   { return AsObject(v) }
func (v Value) AsArray() Array     { return AsArray(v) }
func (v Value) As(ptr any) bool    { return converter.ToAny(v, true, ptr) }

func (v Value) AsBoolOr(def bool) bool          { return AsBoolOr(v, def) }
func (v Value) AsStringOr(def string) string    { return AsStringOr(v, def) }
func (v Value) AsIntOr(def int) int             { return AsIntOr[int](v, def) }
func (v Value) AsInt8Or(def int8) int8          { return AsIntOr[int8](v, def) }
func (v Value) AsInt16Or(def int16) int16       { return AsIntOr[int16](v, def) }
func (v Value) AsInt32Or(def int32) int32       { return AsIntOr[int32](v, def) }
func (v Value) AsInt64Or(def int64) int64       { return AsIntOr[int64](v, def) }
func (v Value) AsUintOr(def uint) uint          { return AsUintOr[uint](v, def) }
func (v Value) AsUint8Or(def uint8) uint8       { return AsUintOr[uint8](v, def) }
func (v Value) AsUint16Or(def uint16) uint16    { return AsUintOr[uint16](v, def) }
func (v Value) AsUint32Or(def uint32) uint32    { return AsUintOr[uint32](v, def) }
func (v Value) AsUint64Or(def uint64) uint64    { return AsUintOr[uint64](v, def) }
func (v Value) AsFloat32Or(def float32) float32 { return AsFloatOr[float32](v, def) }
func (v Value) AsFloat64Or(def float64) float64 { return AsFloatOr[float64](v, def) }
func (v Value) AsObjectOr(def Object) Object    { return AsObjectOr(v, def) }
func (v Value) AsArrayOr(def Array) Array       { return AsArrayOr(v, def) }

func (v Value) get(k string) Value {
	return v.AsObjectOr(nil).Get(k)
}

func unbox(v any) any {
	switch v1 := v.(type) {
	case nil:
		return nil
	case Value:
		return v1.v
	case *Value:
		if v1 == nil {
			return nil
		} else {
			return v1.v
		}
	default:
		return v
	}
}
