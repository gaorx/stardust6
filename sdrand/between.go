package sdrand

import (
	"math/rand"
)

// IntBetween 返回 [low, high) 之间的随机整数
func IntBetween(low, high int) int {
	if low == high {
		return low
	}
	if high < low {
		high, low = low, high
	}
	// [low, high)
	return low + rand.Intn(high-low)
}

// Int64Between 返回 [low, high) 之间的随机整数
func Int64Between(low, high int64) int64 {
	if low == high {
		return low
	}
	if high < low {
		high, low = low, high
	}
	// [low, high)
	return low + rand.Int63n(high-low)
}

// Float64Between 返回 [low, high) 之间的随机浮点数
func Float64Between(low, high float64) float64 {
	if low == high {
		return low
	}
	if high < low {
		high, low = low, high
	}
	// [low, high)
	return low + rand.Float64()*(high-low)
}
