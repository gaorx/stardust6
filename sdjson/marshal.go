package sdjson

import (
	"encoding/json"
)

// UnmarshalBytes 将字节数组反序列化为指定类型
func UnmarshalBytes(j []byte, v any) error {
	return json.Unmarshal(j, v)
}

// MarshalBytes 序列化一个值到JSON形式的字节数组
func MarshalBytes(v any) ([]byte, error) {
	return json.Marshal(v)
}

// MarshalIndentBytes 序列化一个值到JSON形式的字节数组，带缩进
func MarshalIndentBytes(v any, prefix, indent string) ([]byte, error) {
	return json.MarshalIndent(v, prefix, indent)
}

// UnmarshalString 反序列化一个JSON格式字符串到指定类型
func UnmarshalString(s string, v any) error {
	return json.Unmarshal([]byte(s), v)
}

// MarshalString 序列化一个值到JSON格式字符串
func MarshalString(v any) (string, error) {
	j, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(j), nil
}

// MarshalIndentString 序列化一个值到JSON格式字符串，带缩进
func MarshalIndentString(v any, prefix, indent string) (string, error) {
	j, err := json.MarshalIndent(v, prefix, indent)
	if err != nil {
		return "", err
	}
	return string(j), nil
}

// MarshalStringDef 序列化一个值到JSON格式字符串，失败时返回默认值
func MarshalStringDef(v any, def string) string {
	j, err := MarshalString(v)
	if err != nil {
		return def
	}
	return j
}

// MarshalIndentStringDef 序列化一个值到JSON格式字符串，带缩进，失败时返回默认值
func MarshalIndentStringDef(v any, prefix, indent, def string) string {
	j, err := MarshalIndentString(v, prefix, indent)
	if err != nil {
		return def
	}
	return j
}

// MarshalPretty 序列化一个值到JSON格式字符串，带有默认的缩进，失败时返回默认值
func MarshalPretty(v any) string {
	return MarshalIndentStringDef(v, "", "  ", "")
}

// UnmarshalBytesT 将字节数组反序列化为指定类型
func UnmarshalBytesT[T any](j []byte) (T, error) {
	var t T
	if err := json.Unmarshal(j, &t); err != nil {
		var zero T
		return zero, err
	}
	return t, nil
}

// UnmarshalStringT 反序列化一个JSON格式字符串到指定类型
func UnmarshalStringT[T any](j string) (T, error) {
	return UnmarshalBytesT[T]([]byte(j))
}

// UnmarshalBytesDef 将字节数组反序列化为指定类型，失败时返回默认值
func UnmarshalBytesDef[T any](j []byte, def T) T {
	v, err := UnmarshalBytesT[T](j)
	if err != nil {
		return def
	}
	return v
}

// UnmarshalStringDef 反序列化一个JSON格式字符串到指定类型，失败时返回默认值
func UnmarshalStringDef[T any](j string, def T) T {
	return UnmarshalBytesDef[T]([]byte(j), def)
}

// UnmarshalValueBytes 将字节数组反序列化为Value
func UnmarshalValueBytes(j []byte) (Value, error) {
	return UnmarshalBytesT[Value](j)
}

// UnmarshalValueString 反序列化一个JSON格式字符串到Value
func UnmarshalValueString(s string) (Value, error) {
	if v, err := UnmarshalValueBytes([]byte(s)); err != nil {
		return V(nil), err
	} else {
		return v, nil
	}
}

// UnmarshalValueDef 将字节数组反序列化为Value，失败时返回默认值
func UnmarshalValueDef(j []byte, def any) Value {
	v, err := UnmarshalValueBytes(j)
	if err != nil {
		return V(def)
	}
	return v
}

// UnmarshalValueStringDef 反序列化一个JSON格式字符串到Value，失败时返回默认值
func UnmarshalValueStringDef(s string, def any) Value {
	return UnmarshalValueDef([]byte(s), def)
}
