package sdjson

import (
	"github.com/samber/lo"
	"maps"
)

// Object 描述JSON中的Object
type Object map[string]any

// Len 返回对象的键值对数量
func (o Object) Len() int {
	return len(o)
}

// Contains 返回对象是否包含指定键
func (o Object) Contains(k string) bool {
	_, ok := o[k]
	return ok
}

// Get 返回对象中指定键的值
func (o Object) Get(k string) Value {
	v0, ok := o[k]
	if ok {
		return V(v0)
	} else {
		return V(nil)
	}
}

// Set 设置对象中指定键的值
func (o Object) Set(k string, v any) Object {
	if o != nil {
		o[k] = unbox(v)
	}
	return o
}

// First 返回对象中第一个存在的键的值，如果不存在返回nil
func (o Object) First(keys ...string) Value {
	for _, k := range keys {
		v0, ok := o[k]
		if ok {
			return V(v0)
		}
	}
	return V(nil)
}

// Clone 浅拷贝此对象
func (o Object) Clone() Object {
	return maps.Clone(o)
}

// ToPrimitivePossible 尝试将对象中的值转换为基本类型，键不变
func (o Object) ToPrimitivePossible() Object {
	if o == nil {
		return nil
	}
	return lo.MapValues(o, func(v any, _ string) any {
		return ToPrimitivePossible(v)
	})
}
