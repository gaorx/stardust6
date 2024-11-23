package sdparse

import (
	"github.com/gaorx/stardust6/sdjson"
	"time"
)

// Int64 解析到int64，失败返回0
func Int64(s string) int64 {
	return Int64Or(s, 0)
}

// Int 解析到int，失败返回0
func Int(s string) int {
	return IntOr(s, 0)
}

// Uint64 解析到uint64，失败返回0
func Uint64(s string) uint64 {
	return Uint64Or(s, 0)
}

// Uint 解析到uint，失败返回0
func Uint(s string) uint {
	return UintOr(s, 0)
}

// Float64 解析到float64，失败返回0.0
func Float64(s string) float64 {
	return Float64Or(s, 0.0)
}

// Bool 解析到bool，失败返回false
func Bool(s string) bool {
	return BoolOr(s, false)
}

// Time 解析到time.Time，失败返回time.Time{}
func Time(s string) time.Time {
	return TimeOr(s, time.Time{})
}

// JsonValue 解析到sdjson.Value，失败返回sdjson.V(nil)
func JsonValue(s string) sdjson.Value {
	return JsonValueOr(s, sdjson.V(nil))
}

// JsonObject 解析到sdjson.Object，失败返回nil
func JsonObject(s string) sdjson.Object {
	return JsonObjectOr(s, nil)
}

// JsonArray 解析到sdjson.Array，失败返回nil
func JsonArray(s string) sdjson.Array {
	return JsonArrayOr(s, nil)
}
