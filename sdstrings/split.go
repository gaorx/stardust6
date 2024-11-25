package sdstrings

import (
	"strings"
)

// SplitNonempty 将字符串按照指定的分隔符分割
// 如果trimSpace为true，则去掉每个分割后的字符串两边的空白字符，如果分割后的字符串为空，则不加入结果
func SplitNonempty(s, sep string, trimSpace bool) []string {
	r0 := strings.Split(s, sep)
	r1 := make([]string, 0, len(r0))
	for _, a := range r0 {
		if trimSpace {
			a = strings.TrimSpace(a)
		}
		if a != "" {
			r1 = append(r1, a)
		}
	}
	return r1
}

// Split2s 将字符串按照指定的分隔符分割为两部分
func Split2s(s, sep string) (string, string) {
	if s == "" {
		return "", ""
	}
	l := strings.SplitN(s, sep, 2)
	switch len(l) {
	case 0:
		return "", ""
	case 1:
		return l[0], ""
	default:
		return l[0], l[1]
	}
}

// Split3s 将字符串按照指定的分隔符分割为三部分
func Split3s(s, sep string) (string, string, string) {
	if s == "" {
		return "", "", ""
	}
	l := strings.SplitN(s, sep, 3)
	switch len(l) {
	case 0:
		return "", "", ""
	case 1:
		return l[0], "", ""
	case 2:
		return l[0], l[1], ""
	default:
		return l[0], l[1], l[2]
	}
}

// Split4s 将字符串按照指定的分隔符分割为四部分
func Split4s(s, sep string) (string, string, string, string) {
	if s == "" {
		return "", "", "", ""
	}
	l := strings.SplitN(s, sep, 4)
	switch len(l) {
	case 0:
		return "", "", "", ""
	case 1:
		return l[0], "", "", ""
	case 2:
		return l[0], l[1], "", ""
	case 3:
		return l[0], l[1], l[2], ""
	default:
		return l[0], l[1], l[2], l[3]
	}
}
