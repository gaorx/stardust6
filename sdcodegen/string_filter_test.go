package sdcodegen

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStringFilters(t *testing.T) {
	is := assert.New(t)

	// Not
	is.False(HasPrefix("xy", "he").Not()("hello"))

	// HasPrefix
	is.True(HasPrefix("xy", "he")("hello"))
	is.True(HasPrefix("")("hello"))

	// StringEquals
	is.True(StringEquals("hello")("hello"))
	is.False(StringEquals("hell")("hello"))
	is.True(StringEquals("")(""))

	// StringContains
	is.True(StringContains("ell")("hello"))
	is.True(StringContains("")("hello"))

	// HasSuffix
	is.True(HasSuffix("xy", "llo")("hello"))
	is.True(HasSuffix("")("hello"))

	// FilenameMatches
	is.True(FilenameMatches("*.txt", "*.go")("hello.go"))
	is.False(FilenameMatches("*.txt", "*.md")("hello.go"))
	is.False(FilenameMatches("")("hello.go"))

	// RegexpMatches
	is.True(RegexpMatches("^xy$", "^h[a-z]+o$")("hello"))
	is.False(RegexpMatches("^xy$", "^h[0-9]+o$")("hello"))
	is.False(RegexpMatches("")("hello"))
}
