package sdurl

import (
	"github.com/gaorx/stardust6/sdstrings"
	"github.com/samber/lo"
	"strings"
)

// JoinPath 拼接路径，补全其中的'/'
func JoinPath(paths ...string) string {
	var segments []string
	for _, p := range paths {
		segments0 := sdstrings.SplitNonempty(p, "/", false)
		segments = append(segments, segments0...)
	}
	return "/" + strings.Join(lo.Map(segments, func(s string, _ int) string {
		return QueryEscape(s, EscapeEncodePath)
	}), "/")
}
