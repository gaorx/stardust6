package sdauthn

import (
	"time"
)

type Expiration func(start, now time.Time) bool

func (e Expiration) IsExpired(start, now time.Time) bool {
	if e == nil {
		return false
	}
	return e(start, now)
}

func NoExpiration() Expiration {
	return func(_, _ time.Time) bool {
		return false
	}
}

func ExpireAt(expiry time.Time) Expiration {
	return func(_, now time.Time) bool {
		return now.After(expiry)
	}
}

func ExpireIn(duration time.Duration) Expiration {
	return func(start, now time.Time) bool {
		return now.After(start.Add(duration))
	}
}

func ExpireInDays(days int) Expiration {
	return ExpireIn(time.Hour * 24 * time.Duration(days))
}

func ExpireInMinutes(minutes int) Expiration {
	return ExpireIn(time.Minute * time.Duration(minutes))
}
