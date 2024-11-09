package sdjson

import (
	"github.com/samber/lo"
	"slices"
)

// Array 描述JSON中的Array
type Array []any

// Len 返回数组长度
func (a Array) Len() int {
	return len(a)
}

// At 返回数组中的元素，如果越界返回nil
func (a Array) At(i int) Value {
	if i < 0 || i >= len(a) {
		return V(nil)
	}
	return V(a[i])
}

// Set 设置数组中的元素，如果越界不做任何操作
func (a Array) Set(i int, v any) Array {
	if i >= 0 && i < len(a) {
		a[i] = v
	}
	return a
}

// Clone 浅拷贝此数组
func (a Array) Clone() Array {
	return slices.Clone(a)
}

// AsBoolSlice 转换为bool数组
func (a Array) AsBoolSlice() []bool {
	return lo.Map(a, func(v any, _ int) bool {
		return AsBool(v)
	})
}

// AsStringSlice 转换为字符串数组
func (a Array) AsStringSlice() []string {
	return lo.Map(a, func(v any, _ int) string {
		return AsString(v)
	})
}

// AsIntSlice 转换为int数组
func (a Array) AsIntSlice() []int {
	return lo.Map(a, func(v any, _ int) int {
		return AsInt[int](v)
	})
}

// AsInt64Slice 转换为int64数组
func (a Array) AsInt64Slice() []int64 {
	return lo.Map(a, func(v any, _ int) int64 {
		return AsInt[int64](v)
	})
}

// AsFloat64Slice 转换为float64数组
func (a Array) AsFloat64Slice() []float64 {
	return lo.Map(a, func(v any, _ int) float64 {
		return AsFloat[float64](v)
	})
}
