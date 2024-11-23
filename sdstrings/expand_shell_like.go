package sdstrings

import (
	"os"
)

// ExpandMapper 展开字符串值时的映射器
type ExpandMapper = func(string) string

// ExpandMap 生成一个ExpandMapper，将map中的key映射为对应的值
func ExpandMap(m map[string]string) func(string) string {
	return func(k string) string {
		return m[k]
	}
}

// ExpandShellLike 展开一个类似shell环境变量展开的字符串，mapper是将某一个变量名映射为对应的值的函数
func ExpandShellLike(s string, mappers ...ExpandMapper) string {
	return os.Expand(s, func(key string) string {
		for i := len(mappers) - 1; i >= 0; i-- {
			mapper := mappers[i]
			if mapper != nil {
				v := mapper(key)
				if v != "" {
					return v
				}
			}
		}
		return ""
	})
}

// ExpandShellLikeVars 展开一个类似shell环境变量展开的字符串，vars是变量名到值的映射
func ExpandShellLikeVars(s string, vars map[string]string) string {
	return ExpandShellLike(s, ExpandMap(vars))
}
