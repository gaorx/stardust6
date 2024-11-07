package sdreflect

import (
	"reflect"
)

// IsSliceLikeType 判断是否是切片或数组类型
func IsSliceLikeType(t reflect.Type) bool {
	return t.Kind() == reflect.Slice || t.Kind() == reflect.Array
}

// IsSliceLikeValue 判断是否是切片或数组值
func IsSliceLikeValue(v reflect.Value) bool {
	return v.Kind() == reflect.Slice || v.Kind() == reflect.Array
}

// ForEach 遍历切片或数组
func ForEach(v reflect.Value, action func(elem reflect.Value, i, n int)) {
	n := v.Len()
	for i := 0; i < n; i++ {
		action(v.Index(i), i, n)
	}
}
