package sdcodegen

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerate(t *testing.T) {
	is := assert.New(t)
	f, err := Generate("src/a.txt", nil, Text("hello"))
	is.NoError(err)
	is.Equal("src/a.txt", f.Name)
	is.Equal([]byte("hello"), f.Data)
}
