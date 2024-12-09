package sdcheck

import (
	"regexp"
	"strings"
)

// MatchRegexp 生成一个检测函数，检查字符串是否匹配正则表达式
func MatchRegexp(s, pattern string, message any) Func {
	return func() error {
		if matched, err := regexp.MatchString(pattern, s); err != nil {
			return errorOf(message)
		} else {
			if matched {
				return nil
			} else {
				return errorOf(message)
			}
		}
	}
}

// MatchRegexpPattern 生成一个检测函数，检查字符串是否匹配正则表达式
func MatchRegexpPattern(s string, pattern *regexp.Regexp, message any) Func {
	return func() error {
		if matched := pattern.MatchString(s); matched {
			return nil
		} else {
			return errorOf(message)
		}
	}
}

// HasSub 生成一个检测函数，检查字符串是否包含子串
func HasSub(s string, substr string, message any) Func {
	return func() error {
		if !strings.Contains(s, substr) {
			return errorOf(message)
		}
		return nil
	}
}

// HasPrefix 生成一个检测函数，检查字符串是否以前缀开头
func HasPrefix(s string, prefix string, message any) Func {
	return func() error {
		if !strings.HasPrefix(s, prefix) {
			return errorOf(message)
		}
		return nil
	}
}

// HasSuffix 生成一个检测函数，检查字符串是否以后缀结尾
func HasSuffix(s string, suffix string, message any) Func {
	return func() error {
		if !strings.HasSuffix(s, suffix) {
			return errorOf(message)
		}
		return nil
	}
}
