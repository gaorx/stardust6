package sdmath

import (
	"github.com/gaorx/stardust6/sderr"
)

// Interval 区间，表示[min, max)
type Interval[T float64 | float32] struct {
	Min T
	Max T
}

func Normalize[T float64 | float32](v T, src, dst Interval[T]) T {
	if src.Max == src.Min {
		panic(sderr.With("min", src.Min).With("max", src.Max).Newf("illegal source interval"))
	}
	// 归一化, 将value从[src.Min, src.Max]区间映射到[dst.Min, dst.Max]区间,不做参数检查
	return (v-src.Min)/(src.Max-src.Min)*(dst.Max-dst.Min) + dst.Min
}
