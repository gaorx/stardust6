package sdhttpfile

import (
	"github.com/gaorx/stardust6/sdfile"
	"github.com/stretchr/testify/assert"
	"net/http"
	"path/filepath"
	"testing"
)

func TestRead(t *testing.T) {
	is := assert.New(t)
	_ = sdfile.UseTempDir("", "", func(dirname string) {
		file1 := filepath.Join(dirname, "a.txt")
		err := sdfile.WriteText(file1, "hello", 0600)
		is.NoError(err)
		is.FileExists(file1)
		hfs := http.Dir(dirname)
		s, err := ReadText(hfs, "a.txt")
		is.NoError(err)
		is.Equal("hello", s)
		is.Equal("hello", ReadTextDef(hfs, "a.txt", "world"))
		is.Equal("world", ReadTextDef(hfs, "b.txt", "world"))
	})
}
