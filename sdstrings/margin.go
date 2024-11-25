package sdstrings

import (
	"strings"
	"unicode"
)

// TrimMargin 去掉字符串的每行左侧空白和marginPrefix前缀
func TrimMargin(s string, marginPrefix string) string {
	lines := strings.Split(s, "\n")
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		line = strings.TrimLeftFunc(line, func(c rune) bool {
			return unicode.IsSpace(c)
		})
		line = strings.TrimRightFunc(line, func(c rune) bool {
			return c == '\r'
		})
		line = strings.TrimPrefix(line, marginPrefix)
		lines[i] = line
	}
	beginIndex, endIndex := 0, len(lines)
	for i, line := range lines {
		if strings.TrimSpace(line) != "" {
			beginIndex = i
			break
		}
	}
	for i := len(lines) - 1; i >= 0; i-- {
		if strings.TrimSpace(lines[i]) != "" {
			endIndex = i + 1
			break
		}
	}
	if beginIndex > 0 || endIndex < len(lines)-1 {
		lines = lines[beginIndex:endIndex]
	}
	return strings.Join(lines, "\n")
}
