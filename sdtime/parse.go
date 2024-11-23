package sdtime

import (
	"github.com/gaorx/stardust6/sdparse"
	"time"
)

// Parse 解析日期时间
func Parse(s string) (t time.Time, err error) {
	return sdparse.TimeE(s)
}

// ParseOr 解析日期时间，失败时返回默认值
func ParseOr(s string, def time.Time) time.Time {
	return sdparse.TimeOr(s, def)
}

// MustParse 解析日期时间，失败时 panic
func MustParse(s string) time.Time {
	t, err := Parse(s)
	if err != nil {
		panic(err)
	}
	return t
}
