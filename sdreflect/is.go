package sdreflect

import (
	"reflect"
)

// IsStruct 判断是否是结构体值
func IsStruct(t reflect.Type) bool {
	return t.Kind() == reflect.Struct
}

// IsStructPtr 判断是否是结构体指针值
func IsStructPtr(t reflect.Type) bool {
	return t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct
}

// IsSlice 判断是否是切片值，并且元素类型为elemType，如果elemType为nil则不检测元素类型
func IsSlice(t reflect.Type, elemType reflect.Type) bool {
	if t.Kind() != reflect.Slice {
		return false
	}
	if elemType != nil {
		if t.Elem() != elemType {
			return false
		}
	}
	return true
}

// IsSliceLike 判断是否是切片或数组值，并且元素类型为elemType，如果elemType为nil则不检测元素类型
func IsSliceLike(t reflect.Type, elemType reflect.Type) bool {
	if t.Kind() != reflect.Slice && t.Kind() != reflect.Array {
		return false
	}
	if elemType != nil {
		if t.Elem() != elemType {
			return false
		}
	}
	return true
}

// IsMap 判断是否是指定键值类型的映射值，并且keyType和valueType是否为特定类型，
// 如果keyType或valueType为nil则不检测键值类型
func IsMap(t reflect.Type, keyType, valueType reflect.Type) bool {
	if t.Kind() != reflect.Map {
		return false
	}
	if keyType != nil {
		if t.Key() != keyType {
			return false
		}
	}
	if valueType != nil {
		if t.Elem() != valueType {
			return false
		}
	}
	return true
}
