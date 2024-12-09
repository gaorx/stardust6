package sdcheck

import (
	"github.com/gaorx/stardust6/sdvalidate"
)

// ValidateStruct 验证结构体
func ValidateStruct(v any, message string) Func {
	return func() error {
		err := sdvalidate.Struct(v)
		if err != nil {
			return errorOf(message)
		}
		return nil
	}
}

// ValidateStructPartial 验证结构体部分字段
func ValidateStructPartial(v any, fields []string, message string) Func {
	return func() error {
		err := sdvalidate.StructPartial(v, fields)
		if err != nil {
			return errorOf(message)
		}
		return nil
	}
}

// ValidateVar 验证变量
func ValidateVar(v any, tag string, message string) Func {
	return func() error {
		err := sdvalidate.Var(v, tag)
		if err != nil {
			return errorOf(message)
		}
		return nil
	}
}
