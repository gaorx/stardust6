package sdrand

import (
	"github.com/samber/lo"
	"slices"
)

// Shuffle 随机打乱集合
func Shuffle[T any](collection []T) {
	lo.Shuffle(collection)
}

// ShuffleClone 复制集合并随机打乱被复制的集合
func ShuffleClone[T any](collection []T) []T {
	return lo.Shuffle(slices.Clone(collection))
}
