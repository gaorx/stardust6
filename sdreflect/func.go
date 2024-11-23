package sdreflect

import (
	"reflect"
	"slices"
)

// Ins 返回一个函数类型所有的参数类型
func Ins(t reflect.Type) []reflect.Type {
	numIn := t.NumIn()
	ins := make([]reflect.Type, numIn)
	for i := 0; i < numIn; i++ {
		ins[i] = t.In(i)
	}
	return ins
}

// Outs 返回一个函数类型所有的返回值类型
func Outs(t reflect.Type) []reflect.Type {
	numOut := t.NumOut()
	outs := make([]reflect.Type, numOut)
	for i := 0; i < numOut; i++ {
		outs[i] = t.Out(i)
	}
	return outs
}

// TrimLastErr 对一个函数的返回类型进行处理，删除掉最后一个error的类型，返回true则表示原来的返回值中有error类型
func TrimLastErr(outs []reflect.Type) ([]reflect.Type, bool) {
	n := len(outs)
	if n <= 0 {
		return nil, false
	}
	last := outs[n-1]
	if last == TErr {
		return slices.Clone(outs[0 : n-1]), true
	} else {
		return slices.Clone(outs), false
	}
}

// MakeInValues 生成一个函数的参数值
func MakeInValues(inTypes []reflect.Type, f func(t reflect.Type, i int) reflect.Value) []reflect.Value {
	var inVals []reflect.Value
	for i, in := range inTypes {
		outVal := f(in, i)
		inVals = append(inVals, outVal)
	}
	return inVals
}

// MakeInValuesE 生成一个函数的参数值，如果生成过程出错，返回error
func MakeInValuesE(inTypes []reflect.Type, f func(t reflect.Type, i int) (reflect.Value, error)) ([]reflect.Value, error) {
	var inVals []reflect.Value
	for i, in := range inTypes {
		outVal, err := f(in, i)
		if err != nil {
			return nil, err
		}
		inVals = append(inVals, outVal)
	}
	return inVals, nil
}
