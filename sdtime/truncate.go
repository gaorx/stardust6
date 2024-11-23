package sdtime

import (
	"time"
)

// NowTruncateM 截止到时间到分钟级
func NowTruncateM() time.Time {
	return time.Now().Truncate(time.Minute)
}

// NowTruncateS 截止到时间到秒级
func NowTruncateS() time.Time {
	return time.Now().Truncate(time.Second)
}

// NowTruncateMS 截止到时间到毫秒级
func NowTruncateMS() time.Time {
	return time.Now().Truncate(time.Millisecond)
}
