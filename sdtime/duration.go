package sdtime

import (
	"time"
)

// Milliseconds 将毫秒数转换为时间间隔
func Milliseconds(n int64) time.Duration {
	return time.Duration(n) * time.Millisecond
}

// Seconds 将秒数转换为时间间隔
func Seconds(n int64) time.Duration {
	return time.Duration(n) * time.Second
}

// Minutes 将分钟数转换为时间间隔
func Minutes(n int64) time.Duration {
	return time.Duration(n) * time.Minute
}

// Hours 将小时数转换为时间间隔
func Hours(n int64) time.Duration {
	return time.Duration(n) * time.Hour
}

// ToMillis 将时间间隔转换为毫秒数
func ToMillis(d time.Duration) int64 {
	return d.Nanoseconds() / (1000.0 * 1000.0)
}

// ToMillisF 将时间间隔转换为毫秒数，返回浮点数
func ToMillisF(d time.Duration) float64 {
	return float64(d.Nanoseconds() / (1000.0 * 1000.0))
}
