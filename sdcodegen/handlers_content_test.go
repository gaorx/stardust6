package sdcodegen

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestText(t *testing.T) {
	is := assert.New(t)
	s, err := GenerateText(Text("hello"))
	is.NoError(err)
	is.Equal("hello", s)
}
