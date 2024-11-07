package sdcodegen

import (
	"path/filepath"
	"regexp"
	"strings"
)

// StringFilter 字符串过滤器
type StringFilter func(s string) bool

// Not 取反
func (f StringFilter) Not() StringFilter {
	return func(name string) bool {
		return !f(name)
	}
}

// StringEquals 判断字符串是否完全相等
func StringEquals(s string) StringFilter {
	return func(name string) bool {
		return name == s
	}
}

// StringContains 判断字符串是否包含子串
func StringContains(sub string) StringFilter {
	return func(s string) bool {
		return strings.Contains(s, sub)
	}
}

// HasPrefix 判断字符串是否有前缀
func HasPrefix(prefixes ...string) StringFilter {
	return func(s string) bool {
		for _, prefix := range prefixes {
			if strings.HasPrefix(s, prefix) {
				return true
			}
		}
		return false
	}
}

// HasSuffix 判断字符串是否有后缀
func HasSuffix(suffixes ...string) StringFilter {
	return func(s string) bool {
		for _, suffix := range suffixes {
			if strings.HasSuffix(s, suffix) {
				return true
			}
		}
		return false
	}
}

// FilenameMatches 判断文件名是否匹配模式
func FilenameMatches(patterns ...string) StringFilter {
	return func(name string) bool {
		for _, pattern := range patterns {
			if pattern == "" {
				return false
			}
			basename := filepath.Base(name)
			if ok, err := filepath.Match(pattern, basename); ok && err == nil {
				return true
			}
		}
		return false
	}
}

// RegexpMatches 判断字符串是否匹配正则表达式
func RegexpMatches(patterns ...string) StringFilter {
	return func(name string) bool {
		for _, pattern := range patterns {
			if pattern == "" {
				return false
			}
			if ok, err := regexp.MatchString(pattern, name); ok && err == nil {
				return true
			}
		}
		return false
	}
}
