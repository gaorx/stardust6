package sdtime

import (
	"time"
)

// ToUnixS 将时间转换秒级为UNIX时间戳
func ToUnixS(t time.Time) int64 {
	if t.IsZero() {
		return 0
	}
	return t.Unix()
}

// ToUnixMS 将时间转换毫秒级为UNIX时间戳
func ToUnixMS(t time.Time) int64 {
	if t.IsZero() {
		return 0
	}
	return t.UnixNano() / int64(time.Millisecond)
}

// FromUnixS 将秒级的UNIX时间戳转换为时间
func FromUnixS(s int64) time.Time {
	if s == 0 {
		return time.Time{}
	}
	return time.Unix(s, 0)
}

// FromUnixMS 将毫秒级的UNIX时间戳转换为时间
func FromUnixMS(ms int64) time.Time {
	if ms == 0 {
		return time.Time{}
	}
	nanos := ms * 1000000
	return time.Unix(0, nanos)
}

// NowUnixS 获取当前时间的秒级UNIX时间戳
func NowUnixS() int64 {
	return ToUnixS(time.Now())
}

// NowUnixMS 获取当前时间的毫秒级UNIX时间戳
func NowUnixMS() int64 {
	return ToUnixMS(time.Now())
}
