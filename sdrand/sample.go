package sdrand

import (
	"github.com/samber/lo"
)

// Sample 随机从集合中取一个元素
func Sample[T any](collections ...T) T {
	return lo.Sample(collections)
}

// Samples 随机从集合中取 n 个元素
func Samples[T any](collections []T, n int) []T {
	return lo.Samples(collections, n)
}

// W 权重和值的组合
type W[T any] struct {
	W int `json:"w"` // 权重
	V T   `json:"v"` // 值
}

// SampleWeighted 随机从加权集合中取一个元素
func SampleWeighted[T any](collections ...W[T]) T {
	var def T
	n := len(collections)
	if n <= 0 {
		return def
	}
	if n == 1 {
		first := collections[0]
		if first.W > 0 {
			return first.V
		} else {
			return def
		}
	}
	var sum, upto int64 = 0, 0
	for _, w := range collections {
		if w.W > 0 {
			sum += int64(w.W)
		}
	}
	r := Float64Between(0.0, float64(sum))
	for _, w := range collections {
		ww := w.W
		if ww < 0 {
			ww = 0
		}
		if float64(upto)+float64(ww) > r {
			return w.V
		}
		upto += int64(w.W)
	}
	return def
}
