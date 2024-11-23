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

// IsTypes 判断一个类型列表是否是指定的类型列表，如果为都为空则返回true
func IsTypes(actual []reflect.Type, expectant ...reflect.Type) bool {
	if len(actual) != len(expectant) {
		return false
	}
	if len(actual) == 0 {
		return true
	}
	for i, t := range actual {
		if t != expectant[i] {
			return false
		}
	}
	return true
}

// IsAssignableTypes 判断值列表的类型是否是可赋值到指定类型列表，如果为都为空则返回true
func IsAssignableTypes(actual []reflect.Type, expectant ...reflect.Type) bool {
	if len(actual) != len(expectant) {
		return false
	}
	if len(actual) == 0 {
		return true
	}
	for i, t := range actual {
		if !t.AssignableTo(expectant[i]) {
			return false
		}
	}
	return true
}

// IsTypes1 判断类型列表是否是[T1]
func IsTypes1[T1 any](actual []reflect.Type) bool {
	return IsTypes(actual, T[T1]())
}

// IsTypes2 判断类型列表是否是[T1,T2]
func IsTypes2[T1, T2 any](actual []reflect.Type) bool {
	return IsTypes(actual, T[T1](), T[T2]())
}

// IsTypes3 判断类型列表是否是[T1,T2,T3]
func IsTypes3[T1, T2, T3 any](actual []reflect.Type) bool {
	return IsTypes(actual, T[T1](), T[T2](), T[T3]())
}

// IsTypes4 判断类型列表是否是[T1,T2,T3,T4]
func IsTypes4[T1, T2, T3, T4 any](actual []reflect.Type) bool {
	return IsTypes(actual, T[T1](), T[T2](), T[T3](), T[T4]())
}

// IsAssignableTypes1 判断类型列表的类型是否是可赋值到[T1]
func IsAssignableTypes1[T1 any](actual []reflect.Type) bool {
	return IsAssignableTypes(actual, T[T1]())
}

// IsAssignableTypes2 判类型列表的类型是否是可赋值到[T1,T2]
func IsAssignableTypes2[T1, T2 any](actual []reflect.Type) bool {
	return IsAssignableTypes(actual, T[T1](), T[T2]())
}

// IsAssignableTypes3 判断类型列表的类型是否是可赋值到[T1,T2,T3]
func IsAssignableTypes3[T1, T2, T3 any](actual []reflect.Type) bool {
	return IsAssignableTypes(actual, T[T1](), T[T2](), T[T3]())
}

// IsAssignableTypes4 判断类型列表的类型是否是可赋值到[T1,T2,T3,T4]
func IsAssignableTypes4[T1, T2, T3, T4 any](actual []reflect.Type) bool {
	return IsAssignableTypes(actual, T[T1](), T[T2](), T[T3](), T[T4]())
}
