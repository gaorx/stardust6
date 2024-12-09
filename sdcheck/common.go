package sdcheck

import (
	"github.com/gaorx/stardust6/sdreflect"
	"reflect"
)

// True 生成一个检测函数，如果b为false则检测失败
func True(b bool, message any) Func {
	return func() error {
		if !b {
			return errorOf(message)
		}
		return nil
	}
}

// False 生成一个检测函数，如果b为true则检测失败
func False(b bool, message any) Func {
	return func() error {
		if b {
			return errorOf(message)
		}
		return nil
	}
}

// Not 生成一个检测函数，是c的逻辑取反
func Not(c Interface, message any) Func {
	return func() error {
		if c == nil {
			return errorOf(message)
		}
		if err := c.Check(); err == nil {
			return errorOf(message)
		} else {
			return nil
		}
	}
}

// All 将所有检测函数组合成一个检测函数，只要有一个检测失败则整体检测失败
func All(checkers ...Interface) Func {
	if len(checkers) == 0 {
		return Func(nil)
	}
	if len(checkers) == 1 {
		return funcOf(checkers[0])
	}
	return func() error {
		for _, c := range checkers {
			if c != nil {
				if err := c.Check(); err != nil {
					return err
				}
			}
		}
		return nil
	}
}

// And 将所有检测函数组合成一个检测函数，只要有一个检测失败则整体检测失败，并返回错误
func And(checkers []Interface, message any) Func {
	if len(checkers) == 0 {
		return Func(nil)
	}
	if len(checkers) == 1 {
		return funcOf(checkers[0])
	}
	return func() error {
		for _, c := range checkers {
			if c != nil {
				if err := c.Check(); err != nil {
					return errorOf(message)
				}
			}
		}
		return nil
	}
}

// Or 将所有检测函数组合成一个检测函数，只要有一个检测成功则整体检测成功
func Or(checkers []Interface, message any) Func {
	if len(checkers) == 0 {
		return Func(nil)
	}
	if len(checkers) == 1 {
		return funcOf(checkers[0])
	}
	return func() error {
		for _, c := range checkers {
			if c != nil {
				if err := c.Check(); err == nil {
					return nil
				}
			}
		}
		return errorOf(message)
	}
}

// Lazy 延迟加载检测函数，只有在检测的时候才进行调用
type Lazy func() Interface

func (l Lazy) Check() error {
	return l().Check()
}

// If 只有在enabled为true的时候才生成一个基于checker的检测函数，否则生成一个总是成功的检测函数
func If(enabled bool, checker Interface) Func {
	if !enabled {
		return Func(nil)
	}
	return funcOf(checker)
}

// FuncFor 一个检测函数，返回一个值和一个错误
type FuncFor[T any] func() (T, error)

// For 通过一个带返回值的检测函数，生成一个不带返回值的检测函数，返回值则被放置在ptr指向中
func For[T any](f FuncFor[T], ptr *T) Func {
	return func() error {
		r, err := f()
		if err != nil {
			return err
		}
		if ptr != nil {
			*ptr = r
		}
		return nil
	}
}

// Required 生成一个检测函数，如果v为nil或者长度为0则检测失败
func Required(v any, message any) Func {
	return func() error {
		v := sdreflect.RootValueOf(v)
		k := v.Kind()
		if !v.IsValid() {
			return errorOf(message)
		}
		if (k == reflect.Pointer || k == reflect.Func) && v.IsNil() {
			return errorOf(message)
		}
		if (k == reflect.Slice || k == reflect.Array || k == reflect.Map) && (v.IsNil() || v.Len() <= 0) {
			return errorOf(message)
		}
		if v.IsZero() {
			return errorOf(message)
		}
		return nil
	}
}

// Len 生成一个检测函数，检测v的长度是否在[minLen, maxLen]之间
func Len(v any, minLen, maxLen int, message any) Func {
	if maxLen < minLen {
		minLen, maxLen = maxLen, minLen
	}
	return func() error {
		if n := sdreflect.RootValueOf(v).Len(); n < minLen || n > maxLen {
			return errorOf(message)
		}
		return nil
	}
}
