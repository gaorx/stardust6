package sdcodegen

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestResetNewline(t *testing.T) {
	is := assert.New(t)

	s, err := GenerateText(func(c *Context) {
		c.Line("abc")
		c.Line("def")
	}, SetNewline("\n"))
	is.NoError(err)
	is.Equal("abc\ndef\n", s)

	s, err = GenerateText(func(c *Context) {
		c.Line("abc")
		c.Line("def")
	}, SetNewline("\r\n"))
	is.NoError(err)
	is.Equal("abc\r\ndef\r\n", s)

	s, err = GenerateText(func(c *Context) {
		c.Line("abc")
		c.Line("def")
	}, SetNewline("\r\n"), ResetNewline("\n"))
	is.NoError(err)
	is.Equal("abc\ndef\n", s)

	s, err = GenerateText(func(c *Context) {
		c.Line("abc")
		c.Line("def")
	}, SetNewline("\n"), ResetNewline("\r\n"))
	is.NoError(err)
	is.Equal("abc\r\ndef\r\n", s)
}

func TestFinalNewline(t *testing.T) {
	is := assert.New(t)

	s, err := GenerateText(Text("hello"))
	is.NoError(err)
	is.Equal("hello", s)

	s, err = GenerateText(Text("hello\n"), FinalNewline())
	is.NoError(err)
	is.Equal("hello\n", s)

	s, err = GenerateText(Text("hello"), FinalNewline())
	is.NoError(err)
	is.Equal("hello\n", s)
}
