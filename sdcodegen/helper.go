package sdcodegen

import (
	"path/filepath"
	"strings"
)

func toAbs(path string) (string, bool) {
	if filepath.IsAbs(path) {
		return path, true
	}
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", false
	}
	return absPath, true
}

func guessNL(s string) string {
	if strings.Contains(s, "\r\n") {
		return "\r\n"
	}
	return "\n"
}

func mergeMiddlewares(a, b, c []Middleware) []Middleware {
	var r []Middleware
	r = append(r, a...)
	r = append(r, b...)
	r = append(r, c...)
	return r
}
