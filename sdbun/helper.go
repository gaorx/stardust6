package sdbun

import (
	"strings"
)

func ensureRows[T any](rows []T) []T {
	if rows == nil {
		return []T{}
	}
	return rows
}

func trimBrackets(s string) string {
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "(") && strings.HasSuffix(s, ")") {
		s = s[1 : len(s)-1]
	}
	return s
}
