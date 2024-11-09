package sdreflect

import (
	"reflect"
)

// ForEach 遍历切片或数组
func ForEach(v reflect.Value, action func(elem reflect.Value, i, n int)) {
	n := v.Len()
	for i := 0; i < n; i++ {
		action(v.Index(i), i, n)
	}
}
