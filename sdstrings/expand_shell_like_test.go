package sdstrings

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExpandShellLike(t *testing.T) {
	is := assert.New(t)
	is.Equal("abbc", ExpandShellLikeVars("a${b}c", map[string]string{
		"b": "bb",
	}))
	is.Equal("a bb c ddd e ff", ExpandShellLike(
		"a $b c ${d} e $f",
		func(k string) string {
			switch k {
			case "b":
				return "bb"
			case "d":
				return "dd"
			default:
				return ""
			}
		},
		ExpandMap(map[string]string{
			"f": "ff",
			"d": "ddd",
		}),
	))
}
