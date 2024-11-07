package sdreflect

import (
	"github.com/gaorx/stardust6/sderr"
	"reflect"
)

// RootValueOf 返回一个值的reflect.Value，它不会是reflect.Value的嵌套，而是直接包裹值
func RootValueOf(v any) reflect.Value {
	switch v1 := v.(type) {
	case reflect.Value:
		return RootValueOf(v1.Interface())
	case *reflect.Value:
		if v1 == nil {
			panic(sderr.Newf("nil pointer of reflect.Value"))
		}
		return RootValueOf((*v1).Interface())
	default:
		return reflect.ValueOf(v)
	}
}
