package sdcodegen

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddHeaderAndFooter(t *testing.T) {
	is := assert.New(t)
	s, err := GenerateText(func(c *Context) {
		c.Line("abc")
	}, AddHeader("HEADER\n"), AddFooter("FOOTER\n"))
	is.NoError(err)
	is.Equal("HEADER\nabc\nFOOTER\n", s)
}
