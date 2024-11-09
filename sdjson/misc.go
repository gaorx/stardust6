package sdjson

import (
	"encoding/json"
	"github.com/gaorx/stardust6/sderr"
	"github.com/samber/lo"
)

// StructToObject 将一个struct值转换为Object
func StructToObject(v any) (Object, error) {
	j, err := json.Marshal(v)
	if err != nil {
		return nil, sderr.Wrapf(err, "marshal json error")
	}
	var v1 map[string]any
	err = json.Unmarshal(j, &v1)
	if err != nil {
		return nil, sderr.Wrapf(err, "unmarshal to object error")
	}
	return v1, nil
}

// ObjectToStruct 将一个Object转换为struct值
func ObjectToStruct[T any](o Object) (T, error) {
	j, err := json.Marshal(o)
	if err != nil {
		return lo.Empty[T](), sderr.Wrapf(err, "marshal json object error")
	}
	var v1 T
	err = json.Unmarshal(j, &v1)
	if err != nil {
		return lo.Empty[T](), sderr.Wrapf(err, "unmarshal to struct error")
	}
	return v1, nil
}

// ToPrimitive 将一个值转换为基本类型
func ToPrimitive(v any) (any, error) {
	if v == nil {
		return nil, nil
	}
	switch v1 := v.(type) {
	case string, bool:
		return v1, nil
	case int, int8, int16, int32, int64:
		return v1, nil
	case uint, uint8, uint16, uint32, uint64:
		return v1, nil
	case float32, float64:
		return v1, nil
	case json.Number:
		return v1, nil
	default:
		return MarshalString(v)
	}
}

// ToPrimitiveOr 将一个值转换为基本类型，如果转换失败则返回默认值
func ToPrimitiveOr(v any, def any) any {
	v1, err := ToPrimitive(v)
	if err != nil {
		return def
	}
	return v1
}

// ToPrimitivePossible 将一个值转换为基本类型，如果转换失败则返回原值
func ToPrimitivePossible(v any) any {
	v1, err := ToPrimitive(v)
	if err != nil {
		return v
	}
	return v1
}
