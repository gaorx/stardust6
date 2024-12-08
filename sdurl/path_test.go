package sdurl

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJoinPath(t *testing.T) {
	is := assert.New(t)
	is.Equal("/a/b/c/d", JoinPath("a/b", "c/d/", ""))
	is.Equal("/a/%20b/c", JoinPath("a/ b", "c"))
}
