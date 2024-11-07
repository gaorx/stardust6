package sdfile

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRW(t *testing.T) {
	is := assert.New(t)

	var tmpDir string
	err := UseTempDir("", "", func(dirname string) {
		tmpDir = dirname
		filename := filepath.Join(dirname, "a.txt")
		err1 := WriteText(filename, "hello", 0600)
		is.NoError(err1)
		err1 = AppendText(filename, "world", 0600)
		is.NoError(err1)
		s, err1 := ReadText(filename)
		is.NoError(err1)
		is.Equal("helloworld", s)
	})
	is.NoError(err)
	is.NoDirExists(tmpDir)
}
