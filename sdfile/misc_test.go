package sdfile

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestExists(t *testing.T) {
	is := assert.New(t)
	err := UseTempDir("", "", func(dirname string) {
		fn := filepath.Join(dirname, "aaa.txt")
		is.False(Exists(fn))
		err := WriteText(fn, "hello", 0644)
		is.NoError(err)
		is.True(Exists(fn))
		is.False(IsDir(fn))

		dn := filepath.Join(dirname, "subdir")
		is.False(Exists(dn))
		err = os.Mkdir(dn, 0644)
		is.NoError(err)
		is.True(Exists(dn))
		is.True(IsDir(dn))
	})
	is.NoError(err)
}

func TestBinDir(t *testing.T) {
	is := assert.New(t)
	dir := BinDir()
	is.NotEmpty(dir)
	is.DirExists(dir)

	configFn := AbsByBin("a/config.json")
	is.NotEmpty(configFn)
	is.Equal(filepath.Join(dir, "a/config.json"), configFn)
}
