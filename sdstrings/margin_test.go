package sdstrings

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTrimMargin(t *testing.T) {
	is := assert.New(t)
	is.Equal("ABC\n  123\n    456", TrimMargin(`
	|ABC
	|  123
	|    456
	`, "|"))
}
