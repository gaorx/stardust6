package sdstrings

import (
	"github.com/samber/lo"
	"strings"
)

// JoinFunc 将集合中的元素转换为字符串后，使用指定的分隔符连接起来
func JoinFunc[T any](collection []T, sep string, iteratee func(item T, index int) string) string {
	l := lo.Map(collection, iteratee)
	return strings.Join(l, sep)
}
