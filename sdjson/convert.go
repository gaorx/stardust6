package sdjson

import (
	"encoding/json"
	"reflect"
	"strconv"
	"strings"
)

// Converter JSON中的值类型转换器
type Converter interface {
	// ToBool 尝试将一个值转换为bool类型, 如果as为true，则启用跨类型转换
	ToBool(v any, as bool) (bool, bool)

	// ToString 尝试将一个值转换为string类型, 如果as为true，则启用跨类型转换
	ToString(v any, as bool) (string, bool)

	// ToInt 尝试将一个值转换为int类型, 如果as为true，则启用跨类型转换
	ToInt(v any, as bool) (int64, bool)

	// ToUint 尝试将一个值转换为uint类型, 如果as为true，则启用跨类型转换
	ToUint(v any, as bool) (uint64, bool)

	// ToFloat 尝试将一个值转换为float类型, 如果as为true，则启用跨类型转换
	ToFloat(v any, as bool) (float64, bool)

	// ToObject 尝试将一个值转换为Object类型, 如果as为true，则启用跨类型转换
	ToObject(v any, as bool) (Object, bool)

	// ToArray 尝试将一个值转换为Array类型, 如果as为true，则启用跨类型转换
	ToArray(v any, as bool) (Array, bool)

	// ToAny 尝试将一个值转换为任意类型, 如果as为true，则启用跨类型转换, ptr参数转换目标的指针
	ToAny(v any, as bool, ptr any) bool
}

var converter = &mergedConverter{
	converters: []Converter{defaultConverter{}},
}

// Register 注册一个自定义的Converter到类型转换系统中，用于处理自定义类型，例如UUID等
func Register(c Converter) {
	converter.add(c)
}

// ToBool 尝试将一个值转换为bool类型, 如果as为true，则启用跨类型转换
func ToBool(v any, as bool) (bool, bool) {
	return converter.ToBool(v, as)
}

// ToString 尝试将一个值转换为string类型, 如果as为true，则启用跨类型转换
func ToString(v any, as bool) (string, bool) {
	return converter.ToString(v, as)
}

// ToInt 尝试将一个值转换为整数类型, 如果as为true，则启用跨类型转换
func ToInt[T int | int8 | int16 | int32 | int64](v any, as bool) (T, bool) {
	r, ok := converter.ToInt(v, as)
	if !ok {
		return 0, false
	}
	return T(r), true
}

// ToUint 尝试将一个值转换为无符号整数类型, 如果as为true，则启用跨类型转换
func ToUint[T uint | uint8 | uint16 | uint32 | uint64](v any, as bool) (T, bool) {
	r, ok := converter.ToUint(v, as)
	if !ok {
		return 0, false
	}
	return T(r), true
}

// ToFloat 尝试将一个值转换为浮点数类型, 如果as为true，则启用跨类型转换
func ToFloat[T float32 | float64](v any, as bool) (T, bool) {
	r, ok := converter.ToFloat(v, as)
	if !ok {
		return 0.0, false
	}
	return T(r), true
}

// ToObject 尝试将一个值转换为Object类型, 如果as为true，则启用跨类型转换
func ToObject(v any, as bool) (Object, bool) {
	r, ok := converter.ToObject(v, as)
	if !ok {
		return nil, false
	}
	return r, true
}

// ToArray 尝试将一个值转换为Array类型, 如果as为true，则启用跨类型转换
func ToArray(v any, as bool) (Array, bool) {
	r, ok := converter.ToArray(v, as)
	if !ok {
		return nil, false
	}
	return r, true
}

// To 尝试将一个值转换为任意类型, 如果as为true，则启用跨类型转换
func To[T any](v any, as bool) (T, bool) {
	var r T
	ok := converter.ToAny(v, as, &r)
	return r, ok
}

// AsBool 将一个值转换为bool类型，如果转换失败，返回false
func AsBool(v any) bool {
	r, ok := ToBool(v, true)
	if !ok {
		return false
	}
	return r
}

// AsString 将一个值转换为string类型，如果转换失败，返回空字符串
func AsString(v any) string {
	r, ok := ToString(v, true)
	if !ok {
		return ""
	}
	return r
}

// AsInt 将一个值转换为整数类型，如果转换失败，返回0
func AsInt[T int | int8 | int16 | int32 | int64](v any) T {
	r, ok := ToInt[T](v, true)
	if !ok {
		return 0
	}
	return r
}

// AsUint 将一个值转换为无符号整数类型，如果转换失败，返回0
func AsUint[T uint | uint8 | uint16 | uint32 | uint64](v any) T {
	r, ok := ToUint[T](v, true)
	if !ok {
		return 0
	}
	return r
}

// AsFloat 将一个值转换为浮点数类型，如果转换失败，返回0.0
func AsFloat[T float32 | float64](v any) T {
	r, ok := ToFloat[T](v, true)
	if !ok {
		return 0.0
	}
	return r
}

// AsObject 将一个值转换为Object类型，如果转换失败，返回nil
func AsObject(v any) Object {
	r, ok := ToObject(v, true)
	if !ok {
		return nil
	}
	return r
}

// AsArray 将一个值转换为Array类型，如果转换失败，返回nil
func AsArray(v any) Array {
	r, ok := ToArray(v, true)
	if !ok {
		return nil
	}
	return r
}

// As 将一个值转换为任意类型，如果转换失败，返回该类型的零值
func As[T any](v any) T {
	var r T
	ok := converter.ToAny(v, true, &r)
	if !ok {
		var zero T
		return zero
	}
	return r
}

// AsBoolOr 将一个值转换为bool类型，如果转换失败，返回默认值
func AsBoolOr(v any, def bool) bool {
	r, ok := ToBool(v, true)
	if !ok {
		return def
	}
	return r
}

// AsStringOr 将一个值转换为string类型，如果转换失败，返回默认值
func AsStringOr(v any, def string) string {
	r, ok := ToString(v, true)
	if !ok {
		return def
	}
	return r
}

// AsIntOr 将一个值转换为整数类型，如果转换失败，返回默认值
func AsIntOr[T int | int8 | int16 | int32 | int64](v any, def T) T {
	r, ok := ToInt[T](v, true)
	if !ok {
		return def
	}
	return r
}

// AsUintOr 将一个值转换为无符号整数类型，如果转换失败，返回默认值
func AsUintOr[T uint | uint8 | uint16 | uint32 | uint64](v any, def T) T {
	r, ok := ToUint[T](v, true)
	if !ok {
		return def
	}
	return r
}

// AsFloatOr 将一个值转换为浮点数类型，如果转换失败，返回默认值
func AsFloatOr[T float32 | float64](v any, def T) T {
	r, ok := ToFloat[T](v, true)
	if !ok {
		return def
	}
	return r
}

// AsObjectOr 将一个值转换为Object类型，如果转换失败，返回默认值
func AsObjectOr(v any, def Object) Object {
	r, ok := ToObject(v, true)
	if !ok {
		return def
	}
	return r
}

// AsArrayOr 将一个值转换为Array类型，如果转换失败，返回默认值
func AsArrayOr(v any, def Array) Array {
	r, ok := ToArray(v, true)
	if !ok {
		return def
	}
	return r
}

// AsOr 将一个值转换为任意类型，如果转换失败，返回默认值
func AsOr[T any](v any, def T) T {
	r, ok := To[T](v, true)
	if !ok {
		return def
	}
	return r
}

// merged converter

type mergedConverter struct {
	converters []Converter
}

func (mc *mergedConverter) add(c Converter) {
	if c != nil {
		mc.converters = append(mc.converters, c)
	}
}

func (mc *mergedConverter) ToBool(v any, as bool) (bool, bool) {
	for _, c := range mc.converters {
		v, ok := c.ToBool(v, as)
		if ok {
			return v, ok
		}
	}
	return false, false
}

func (mc *mergedConverter) ToString(v any, as bool) (string, bool) {
	for _, c := range mc.converters {
		v, ok := c.ToString(v, as)
		if ok {
			return v, ok
		}
	}
	return "", false
}

func (mc *mergedConverter) ToInt(v any, as bool) (int64, bool) {
	for _, c := range mc.converters {
		v, ok := c.ToInt(v, as)
		if ok {
			return v, ok
		}
	}
	return 0, false
}

func (mc *mergedConverter) ToUint(v any, as bool) (uint64, bool) {
	for _, c := range mc.converters {
		v, ok := c.ToUint(v, as)
		if ok {
			return v, ok
		}
	}
	return 0, false
}

func (mc *mergedConverter) ToFloat(v any, as bool) (float64, bool) {
	for _, c := range mc.converters {
		v, ok := c.ToFloat(v, as)
		if ok {
			return v, ok
		}
	}
	return 0, false
}

func (mc *mergedConverter) ToObject(v any, as bool) (Object, bool) {
	for _, c := range mc.converters {
		v, ok := c.ToObject(v, as)
		if ok {
			return v, ok
		}
	}
	return nil, false
}

func (mc *mergedConverter) ToArray(v any, as bool) (Array, bool) {
	for _, c := range mc.converters {
		v, ok := c.ToArray(v, as)
		if ok {
			return v, ok
		}
	}
	return nil, false
}

func (mc *mergedConverter) ToAny(v any, as bool, ptr any) bool {
	for _, c := range mc.converters {
		ok := c.ToAny(v, as, ptr)
		if ok {
			return true
		}
	}
	return false
}

// default converter

type defaultConverter struct{}

func (c defaultConverter) ToBool(v any, as bool) (bool, bool) {
	v = unbox(v)

	// nil
	if v == nil {
		return false, false
	}

	// bool
	if v1, ok := v.(bool); ok {
		return v1, true
	}
	if as {
		switch v1 := v.(type) {
		// string
		case string:
			b, err := strconv.ParseBool(v1)
			if err != nil {
				return false, false
			}
			return b, true

		// number
		case json.Number:
			if isFloat(v1) {
				if v2, err := v1.Float64(); err != nil {
					return false, false
				} else {
					return v2 != 0.0, true
				}
			} else {
				if v2, err := v1.Int64(); err != nil {
					return false, false
				} else {
					return v2 != 0, true
				}
			}
		case int:
			return v1 != 0, true
		case int64:
			return v1 != 0, true
		case uint:
			return v1 != 0, true
		case uint64:
			return v1 != 0, true
		case float64:
			return v1 != 0.0, true
		case int8:
			return v1 != 0, true
		case int16:
			return v1 != 0, true
		case int32:
			return v1 != 0, true
		case uint8:
			return v1 != 0, true
		case uint16:
			return v1 != 0, true
		case uint32:
			return v1 != 0, true
		case float32:
			return v1 != 0.0, true
		}
	}
	return false, false
}

func (c defaultConverter) ToString(v any, as bool) (string, bool) {
	v = unbox(v)

	// nil
	if v == nil {
		return "", false
	}
	// string
	if v1, ok := v.(string); ok {
		return v1, true
	}

	if as {
		switch v1 := v.(type) {
		// bool
		case bool:
			return strconv.FormatBool(v1), true
		// number
		case json.Number:
			return v1.String(), true
		case int:
			return strconv.FormatInt(int64(v1), 10), true
		case int64:
			return strconv.FormatInt(v1, 10), true
		case uint:
			return strconv.FormatUint(uint64(v1), 10), true
		case uint64:
			return strconv.FormatUint(v1, 10), true
		case float64:
			return strconv.FormatFloat(v1, 'f', -1, 64), true
		case int8:
			return strconv.FormatInt(int64(v1), 10), true
		case int16:
			return strconv.FormatInt(int64(v1), 10), true
		case int32:
			return strconv.FormatInt(int64(v1), 10), true
		case uint8:
			return strconv.FormatUint(uint64(v1), 10), true
		case uint16:
			return strconv.FormatUint(uint64(v1), 10), true
		case uint32:
			return strconv.FormatUint(uint64(v1), 10), true
		case float32:
			return strconv.FormatFloat(float64(v1), 'f', -1, 32), true
		}
	}

	return "", false
}

func (c defaultConverter) ToInt(v any, as bool) (int64, bool) {
	v = unbox(v)

	// nil
	if v == nil {
		return 0, false
	}

	// number
	switch v1 := v.(type) {
	case json.Number:
		if isFloat(v1) {
			if v2, err := v1.Float64(); err != nil {
				return 0, false
			} else {
				return int64(v2), true
			}
		} else {
			if v2, err := v1.Int64(); err != nil {
				return 0, false
			} else {
				return v2, true
			}
		}
	case int:
		return int64(v1), true
	case int64:
		return v1, true
	case uint:
		return int64(v1), true
	case uint64:
		return int64(v1), true
	case float64:
		return int64(v1), true
	case int8:
		return int64(v1), true
	case int16:
		return int64(v1), true
	case int32:
		return int64(v1), true
	case uint8:
		return int64(v1), true
	case uint16:
		return int64(v1), true
	case uint32:
		return int64(v1), true
	case float32:
		return int64(v1), true
	}

	if as {
		switch v1 := v.(type) {
		// bool
		case bool:
			if v1 {
				return 1, true
			} else {
				return 0, true
			}
		// string
		case string:
			if isFloat(json.Number(v1)) {
				if v2, err := json.Number(v1).Float64(); err != nil {
					return 0, false
				} else {
					return int64(v2), true
				}
			} else {
				if v2, err := json.Number(v1).Int64(); err != nil {
					return 0, false
				} else {
					return v2, true
				}
			}
		}
	}

	return 0, false
}

func (c defaultConverter) ToUint(v any, as bool) (uint64, bool) {
	v = unbox(v)

	// nil
	if v == nil {
		return 0, false
	}

	// number
	switch v1 := v.(type) {
	case json.Number:
		if isFloat(v1) {
			if v2, err := v1.Float64(); err != nil {
				return 0, false
			} else {
				return uint64(v2), true
			}
		} else {
			if v2, err := strconv.ParseUint(string(v1), 10, 64); err != nil {
				return 0, false
			} else {
				return v2, true
			}
		}
	case int:
		return uint64(v1), true
	case int64:
		return uint64(v1), true
	case uint:
		return uint64(v1), true
	case uint64:
		return v1, true
	case float64:
		return uint64(v1), true
	case int8:
		return uint64(v1), true
	case int16:
		return uint64(v1), true
	case int32:
		return uint64(v1), true
	case uint8:
		return uint64(v1), true
	case uint16:
		return uint64(v1), true
	case uint32:
		return uint64(v1), true
	case float32:
		return uint64(v1), true
	}

	if as {
		switch v1 := v.(type) {
		// bool
		case bool:
			if v1 {
				return 1, true
			} else {
				return 0, true
			}
		// string
		case string:
			if isFloat(json.Number(v1)) {
				if v2, err := json.Number(v1).Float64(); err != nil {
					return 0, false
				} else {
					return uint64(v2), true
				}
			} else {
				if v2, err := strconv.ParseUint(v1, 10, 64); err != nil {
					return 0, false
				} else {
					return v2, true
				}
			}
		}
	}

	return 0, false
}

func (c defaultConverter) ToFloat(v any, as bool) (float64, bool) {
	v = unbox(v)

	// nil
	if v == nil {
		return 0, false
	}

	// number
	switch v1 := v.(type) {
	case json.Number:
		if v2, err := v1.Float64(); err != nil {
			return 0, false
		} else {
			return v2, true
		}
	case float64:
		return v1, true
	case float32:
		return float64(v1), true
	case int:
		return float64(v1), true
	case int64:
		return float64(v1), true
	case uint:
		return float64(v1), true
	case uint64:
		return float64(v1), true
	case int8:
		return float64(v1), true
	case int16:
		return float64(v1), true
	case int32:
		return float64(v1), true
	case uint8:
		return float64(v1), true
	case uint16:
		return float64(v1), true
	case uint32:
		return float64(v1), true
	}

	if as {
		switch v1 := v.(type) {
		// bool
		case bool:
			if v1 {
				return 1.0, true
			} else {
				return 0.0, true
			}
		// string
		case string:
			if v2, err := json.Number(v1).Float64(); err != nil {
				return 0.0, false
			} else {
				return v2, true
			}
		}
	}

	return 0.0, false
}

func (c defaultConverter) ToObject(v any, as bool) (Object, bool) {
	v = unbox(v)

	// nil
	if v == nil {
		return nil, false
	}

	// map like
	if v1, ok := v.(map[string]any); ok {
		return v1, true
	} else if v1, ok := v.(Object); ok {
		return v1, true
	} else if rv := reflect.ValueOf(v); rv.Type().Kind() == reflect.Map && rv.Type().Key().Kind() == reflect.String {
		return genericMapToObject(rv)
	}

	if as {
		// struct
		rv := reflect.ValueOf(v)
		rt := rv.Type()
		if rt.Kind() == reflect.Struct {
			return structToObject(rv)
		} else if rt.Kind() == reflect.Ptr && rt.Elem().Kind() == reflect.Struct {
			return structToObject(rv)
		}
	}

	return nil, false
}

func (c defaultConverter) ToArray(v any, as bool) (Array, bool) {
	v = unbox(v)

	// nil
	if v == nil {
		return nil, false
	}

	// slice like
	if v1, ok := v.([]any); ok {
		return v1, true
	} else if v1, ok := v.(Array); ok {
		return v1, true
	} else if rv := reflect.ValueOf(v); rv.Type().Kind() == reflect.Slice || rv.Type().Kind() == reflect.Array {
		return genericSliceToArray(rv)
	}

	return nil, false
}

func (c defaultConverter) ToAny(v any, as bool, ptr any) bool {
	raw, err := json.Marshal(v)
	if err != nil {
		return false
	}
	err = json.Unmarshal(raw, ptr)
	if err != nil {
		return false
	}
	return true
}

func genericMapToObject(rv reflect.Value) (map[string]any, bool) {
	if rv.IsNil() {
		return nil, true
	}
	l := rv.Len()
	m := make(map[string]any, l)
	iter := rv.MapRange()
	for iter.Next() {
		k := iter.Key().Interface().(string)
		v := iter.Value().Interface()
		m[k] = v
	}
	return m, true
}

func structToObject(rv reflect.Value) (map[string]any, bool) {
	if rv.Kind() == reflect.Ptr && rv.IsNil() {
		return nil, true
	}
	raw, err := json.Marshal(rv.Interface())
	if err != nil {
		return nil, false
	}
	var m map[string]any
	err = json.Unmarshal(raw, &m)
	if err != nil {
		return nil, false
	}
	return m, true
}

func genericSliceToArray(rv reflect.Value) ([]any, bool) {
	if rv.Kind() == reflect.Slice && rv.IsNil() {
		return nil, true
	}
	l := rv.Len()
	a := make([]any, 0, l)
	for i := 0; i < l; i++ {
		a = append(a, rv.Index(i).Interface())
	}
	return a, true
}

func isFloat(n json.Number) bool {
	return strings.Contains(n.String(), ".")
}
