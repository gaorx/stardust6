package sdcodegen

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPlaceholder(t *testing.T) {
	is := assert.New(t)
	text, err := GenerateText(func(c *Context) {
		c.Line("HEADER")
		c.Placeholder("content")
		c.Line("----")
		c.Placeholder("content")
		c.Line("----")
		c.Placeholder("no_content").Newl()
		c.Line("FOOTER")
		c.ExpandPlaceholder("content", func() {
			c.Line("XYZ")
		})
	})
	is.NoError(err)
	is.Equal("HEADER\nXYZ\n----\nXYZ\n----\n\nFOOTER\n", text)
}
