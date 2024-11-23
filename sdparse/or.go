package sdparse

import (
	"github.com/gaorx/stardust6/sdjson"
	"time"
)

// Int64Or 解析到int64，失败返回默认值
func Int64Or(s string, def int64) int64 {
	r, err := Int64E(s)
	if err != nil {
		return def
	}
	return r
}

// IntOr 解析到int，失败返回默认值
func IntOr(s string, def int) int {
	r, err := IntE(s)
	if err != nil {
		return def
	}
	return r
}

// Uint64Or 解析到uint64，失败返回默认值
func Uint64Or(s string, def uint64) uint64 {
	r, err := Uint64E(s)
	if err != nil {
		return def
	}
	return r
}

// UintOr 解析到uint，失败返回默认值
func UintOr(s string, def uint) uint {
	r, err := UintE(s)
	if err != nil {
		return def
	}
	return r
}

// Float64Or 解析到float64，失败返回默认值
func Float64Or(s string, def float64) float64 {
	r, err := Float64E(s)
	if err != nil {
		return def
	}
	return r
}

// BoolOr 解析到bool，失败返回默认值
func BoolOr(s string, def bool) bool {
	r, err := BoolE(s)
	if err != nil {
		return def
	}
	return r
}

// TimeOr 解析到time.Time，失败返回默认值
func TimeOr(s string, def time.Time) time.Time {
	r, err := TimeE(s)
	if err != nil {
		return def
	}
	return r
}

// JsonValueOr 解析到sdjson.Value，失败返回默认值
func JsonValueOr(s string, def any) sdjson.Value {
	r, err := JsonValueE(s)
	if err != nil {
		return sdjson.V(def)
	}
	return r
}

// JsonObjectOr 解析到sdjson.Object，失败返回默认值
func JsonObjectOr(s string, def sdjson.Object) sdjson.Object {
	r, err := JsonObjectE(s)
	if err != nil {
		return def
	}
	return r
}

// JsonArrayOr 解析到sdjson.Array，失败返回默认值
func JsonArrayOr(s string, def sdjson.Array) sdjson.Array {
	r, err := JsonArrayE(s)
	if err != nil {
		return def
	}
	return r
}
