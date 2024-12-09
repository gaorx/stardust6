package sdcheck

import (
	"slices"
)

// In 生成一个检测函数，检查v是否在available中
func In[T comparable, C ~[]T](v T, available C, message any) Func {
	return func() error {
		if !slices.Contains[C, T](available, v) {
			return errorOf(message)
		}
		return nil
	}
}

// NotIn 生成一个检测函数，检查v是否不在available中
func NotIn[T comparable, C ~[]T](v T, available C, message any) Func {
	return func() error {
		if slices.Contains[C, T](available, v) {
			return errorOf(message)
		}
		return nil
	}
}

// HasKey 生成一个检测函数，检查m是否包含k
func HasKey[K comparable, V any, M ~map[K]V](k K, m M, message any) Func {
	return func() error {
		if _, ok := m[k]; !ok {
			return errorOf(message)
		}
		return nil
	}
}

// NotHasKey 生成一个检测函数，检查m是否不包含k
func NotHasKey[K comparable, V any, M ~map[K]V](k K, m M, message any) Func {
	return func() error {
		if _, ok := m[k]; ok {
			return errorOf(message)
		}
		return nil
	}
}
