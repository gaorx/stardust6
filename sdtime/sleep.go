package sdtime

import (
	"time"
)

// SleepM 休眠指定分钟数
func SleepM(n int64) {
	time.Sleep(Minutes(n))
}

// SleepS 休眠指定秒数
func SleepS(n int64) {
	time.Sleep(Seconds(n))
}

// SleepMS 休眠指定毫秒数
func SleepMS(n int64) {
	time.Sleep(Milliseconds(n))
}
