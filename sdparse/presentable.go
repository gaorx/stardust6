package sdparse

import (
	"github.com/gaorx/stardust6/sdjson"
	"time"
)

// Presentable 用于可以转换成其他格式的字符串
type Presentable string

// String 返回字符串
func (p Presentable) String() string {
	return string(p)
}

// Int64E 转换为 int64，如果失败返回错误
func (p Presentable) Int64E() (int64, error) {
	return Int64E(string(p))
}

// IntE 转换为 int，如果失败返回错误
func (p Presentable) IntE() (int, error) {
	return IntE(string(p))
}

// Uint64E 转换为 uint64，如果失败返回错误
func (p Presentable) Uint64E() (uint64, error) {
	return Uint64E(string(p))
}

// UintE 转换为 uint，如果失败返回错误
func (p Presentable) UintE() (uint, error) {
	return UintE(string(p))
}

// Float64E 转换为 float64，如果失败返回错误
func (p Presentable) Float64E() (float64, error) {
	return Float64E(string(p))
}

// BoolE 转换为 bool，如果失败返回错误
func (p Presentable) BoolE() (bool, error) {
	return BoolE(string(p))
}

// TimeE 转换为 time.Time，如果失败返回错误
func (p Presentable) TimeE() (time.Time, error) {
	return TimeE(string(p))
}

// JsonValueE 转换为 sdjson.Value，如果失败返回错误
func (p Presentable) JsonValueE() (sdjson.Value, error) {
	return JsonValueE(string(p))
}

// JsonObjectE 转换为 sdjson.Object，如果失败返回错误
func (p Presentable) JsonObjectE() (sdjson.Object, error) {
	return JsonObjectE(string(p))
}

// JsonArrayE 转换为 sdjson.Array，如果失败返回错误
func (p Presentable) JsonArrayE() (sdjson.Array, error) {
	return JsonArrayE(string(p))
}

// Int64 转换为 int64
func (p Presentable) Int64() int64 {
	return Int64(string(p))
}

// Int 转换为 int
func (p Presentable) Int() int {
	return Int(string(p))
}

// Uint64 转换为 uint64
func (p Presentable) Uint64() uint64 {
	return Uint64(string(p))
}

// Uint 转换为 uint
func (p Presentable) Uint() uint {
	return Uint(string(p))
}

// Float64 转换为 float64
func (p Presentable) Float64() float64 {
	return Float64(string(p))
}

// Bool 转换为 bool
func (p Presentable) Bool() bool {
	return Bool(string(p))
}

// Time 转换为 time.Time
func (p Presentable) Time() time.Time {
	return Time(string(p))
}

// JsonValue 转换为 sdjson.Value
func (p Presentable) JsonValue() sdjson.Value {
	return JsonValue(string(p))
}

// JsonObject 转换为 sdjson.Object
func (p Presentable) JsonObject() sdjson.Object {
	return JsonObject(string(p))
}

// JsonArray 转换为 sdjson.Array
func (p Presentable) JsonArray() sdjson.Array {
	return JsonArray(string(p))
}

// Int64Or 转换为 int64，如果失败返回默认值
func (p Presentable) Int64Or(def int64) int64 {
	return Int64Or(string(p), def)
}

// IntOr 转换为 int，如果失败返回默认值
func (p Presentable) IntOr(def int) int {
	return IntOr(string(p), def)
}

// Uint64Or 转换为 uint64，如果失败返回默认值
func (p Presentable) Uint64Or(def uint64) uint64 {
	return Uint64Or(string(p), def)
}

// UintOr 转换为 uint，如果失败返回默认值
func (p Presentable) UintOr(def uint) uint {
	return UintOr(string(p), def)
}

// Float64Or 转换为 float64，如果失败返回默认值
func (p Presentable) Float64Or(def float64) float64 {
	return Float64Or(string(p), def)
}

// BoolOr 转换为 bool，如果失败返回默认值
func (p Presentable) BoolOr(def bool) bool {
	return BoolOr(string(p), def)
}

// TimeOr 转换为 time.Time，如果失败返回默认值
func (p Presentable) TimeOr(def time.Time) time.Time {
	return TimeOr(string(p), def)
}

// JsonValueOr 转换为 sdjson.Value，如果失败返回默认值
func (p Presentable) JsonValueOr(def any) sdjson.Value {
	return JsonValueOr(string(p), def)
}

// JsonObjectOr 转换为 sdjson.Object，如果失败返回默认值
func (p Presentable) JsonObjectOr(def sdjson.Object) sdjson.Object {
	return JsonObjectOr(string(p), def)
}

// JsonArrayOr 转换为 sdjson.Array，如果失败返回默认值
func (p Presentable) JsonArrayOr(def sdjson.Array) sdjson.Array {
	return JsonArrayOr(string(p), def)
}
