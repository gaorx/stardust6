package sdauthn

import (
	"time"
)

// Expiration 过期判断函数；start是开始时间，now是当前时间。
type Expiration func(start, now time.Time) bool

// IsExpired 是否过期
func (e Expiration) IsExpired(start, now time.Time) bool {
	if e == nil {
		return false
	}
	return e(start, now)
}

// NoExpiration 永不过期
func NoExpiration() Expiration {
	return func(_, _ time.Time) bool {
		return false
	}
}

// ExpireAt 在某时刻过期
func ExpireAt(expiry time.Time) Expiration {
	return func(_, now time.Time) bool {
		return now.After(expiry)
	}
}

// ExpireIn 在duration时间后过期
func ExpireIn(duration time.Duration) Expiration {
	return func(start, now time.Time) bool {
		return now.After(start.Add(duration))
	}
}

// ExpireInDays 在days天后过期
func ExpireInDays(days int) Expiration {
	return ExpireIn(time.Hour * 24 * time.Duration(days))
}

// ExpireInMinutes 在minutes分钟后过期
func ExpireInMinutes(minutes int) Expiration {
	return ExpireIn(time.Minute * time.Duration(minutes))
}
