package sderr

import (
	"strconv"
)

func ensurePtr[T any](p *T) *T {
	if p == nil {
		return new(T)
	}
	return p
}

func quote(s string) string {
	return strconv.Quote(s)
}
