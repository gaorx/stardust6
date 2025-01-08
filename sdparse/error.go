package sdparse

import (
	"github.com/gaorx/stardust6/sdjson"
	"strconv"
	"time"

	"github.com/gaorx/stardust6/sderr"
)

// Int64E 解析到int64，失败返回错误
func Int64E(s string) (int64, error) {
	r, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, sderr.Wrap(err)
	}
	return r, nil
}

// IntE 解析到int，失败返回错误
func IntE(s string) (int, error) {
	r, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, sderr.Wrap(err)
	}
	return int(r), nil
}

// Uint64E 解析到uint64，失败返回错误
func Uint64E(s string) (uint64, error) {
	r, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, sderr.Wrap(err)
	}
	return r, nil
}

// UintE 解析到uint，失败返回错误
func UintE(s string) (uint, error) {
	r, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, sderr.Wrap(err)
	}
	return uint(r), nil
}

// Float64E 解析到float64，失败返回错误
func Float64E(s string) (float64, error) {
	r, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0.0, sderr.Wrap(err)
	}
	return r, nil
}

// BoolE 解析到bool，失败返回错误
func BoolE(s string) (bool, error) {
	r, err := strconv.ParseBool(s)
	if err != nil {
		return false, sderr.Wrap(err)
	}
	return r, nil
}

var (
	timeLayoutsForParse = []string{
		"20060102150405",
		time.ANSIC,
		time.UnixDate,
		time.RubyDate,
		time.Kitchen,
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02",                         // RFC 3339
		"2006-01-02 15:04",                   // RFC 3339 with minutes
		"2006-01-02 15:04:05",                // RFC 3339 with seconds
		"2006-01-02 15:04:05-07:00",          // RFC 3339 with seconds and timezone
		"2006-01-02T15Z0700",                 // ISO8601 with hour
		"2006-01-02T15:04Z0700",              // ISO8601 with minutes
		"2006-01-02T15:04:05Z0700",           // ISO8601 with seconds
		"2006-01-02T15:04:05.999999999Z0700", // ISO8601 with nanoseconds
	}
)

// TimeE 解析到time.Time，失败返回错误
func TimeE(s string) (time.Time, error) {
	for _, layout := range timeLayoutsForParse {
		r, err := time.Parse(layout, s)
		if err == nil {
			return r, nil
		}
	}
	return time.Time{}, sderr.Newf("parse time failed")
}

// JsonValueE 解析到sdjson.Value，失败返回错误
func JsonValueE(s string) (sdjson.Value, error) {
	return sdjson.UnmarshalValueString(s)
}

// JsonObjectE 解析到sdjson.Objects，失败返回错误
func JsonObjectE(s string) (sdjson.Object, error) {
	jv, err := JsonValueE(s)
	if err != nil {
		return nil, err
	}
	o, ok := jv.ToObject(true)
	if !ok {
		return nil, sderr.Newf("not a json object")
	}
	return o, nil
}

// JsonArrayE 解析到sdjson.Array，失败返回错误
func JsonArrayE(s string) (sdjson.Array, error) {
	jv, err := JsonValueE(s)
	if err != nil {
		return nil, err
	}
	a, ok := jv.ToArray(true)
	if !ok {
		return nil, sderr.Newf("not a json array")
	}
	return a, nil
}
