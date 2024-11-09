// Package sdjson JSON工具库
//
// 本库提供了一些JSON相关的工具函数，主要功能包括类型转换与序列化。
//
// 1. 类型转换：与JSON定义一致，逻辑类型包括array、object、string、number、boolean、null。 对应的Go类型为：
//   - array = []any
//   - object = map[string]any
//   - string = string
//   - number = int/int8/int16/int32/int64/uint/uint8/uint16/uint32/uint64/float32/float64/json.Number
//   - boolean = bool
//   - null = nil
//
// 所有ToXxx()都有第二个参数as，如果这个参数为true的话，可以在不同的逻辑类型转换，否则只能在相同逻辑类型之间转换。
// 例如int和float64同属number逻辑类型，可以转换，而int和string则不能转换;
// 所有AsXxx()的函数都是在不同逻辑类型之间转换，例如int和string互相转换。
//
// 2. Marshal/Unmarshal 支持了bytes为参数的系列，也支持string为参数的系列函数，也在unmarshal中支持了范型。
package sdjson
