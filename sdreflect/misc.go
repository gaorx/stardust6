package sdreflect

import (
	"github.com/gaorx/stardust6/sderr"
	"github.com/samber/lo"
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

// Deref 返回一个值的非指针值，如果是指针则递归解引用
func Deref(v reflect.Value) reflect.Value {
	if v.Kind() == reflect.Ptr {
		return Deref(v.Elem())
	}
	return v
}

// TypesOf 返回reflect.Value列表的类型列表
func TypesOf(vals []reflect.Value) []reflect.Type {
	return lo.Map(vals, func(v reflect.Value, _ int) reflect.Type {
		return v.Type()
	})
}

// ToTypes 将任意值列表转换为reflect.Type列表
func ToTypes(a ...any) []reflect.Type {
	return lo.Map(a, func(x any, i int) reflect.Type {
		return reflect.TypeOf(x)
	})
}

// ToValues 将任意值列表转换为reflect.Value列表
func ToValues(a ...any) []reflect.Value {
	return lo.Map(a, func(x any, i int) reflect.Value {
		return reflect.ValueOf(x)
	})
}

// ToInterfaces 将reflect.Value列表转换为任意值列表
func ToInterfaces(a ...reflect.Value) []any {
	return lo.Map(a, func(x reflect.Value, i int) any {
		return x.Interface()
	})
}
