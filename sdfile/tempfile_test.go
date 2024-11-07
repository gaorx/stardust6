package sdfile

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTempFile(t *testing.T) {
	is := assert.New(t)
	const text = "hello"
	var tmpName string
	s, err := UseTempFileFor("", "", func(f *os.File) (string, error) {
		tmpName = f.Name()
		_, err := f.WriteString(text)
		is.NoError(err)
		_, err = f.Seek(0, io.SeekStart)
		is.NoError(err)
		d, err := io.ReadAll(f)
		is.NoError(err)
		return string(d), nil
	})
	is.NoError(err)
	is.Equal(text, s)
	is.True(tmpName != "")
	is.NoFileExists(tmpName)

	err = UseTempFile("", "", func(f *os.File) {
		tmpName = f.Name()
		_, err := f.WriteString(text)
		is.NoError(err)
		_, err = f.Seek(0, io.SeekStart)
		is.NoError(err)
		d, err := io.ReadAll(f)
		is.NoError(err)
		is.Equal([]byte(text), d)
	})
	is.NoError(err)
	is.True(tmpName != "")
	is.NoFileExists(tmpName)
}

func TestTempDir(t *testing.T) {
	is := assert.New(t)
	const text = "hello"
	var tmpDir string
	s, err := UseTempDirFor("", "", func(dirname string) (string, error) {
		tmpDir = dirname
		filename := filepath.Join(dirname, "b.txt")
		err1 := WriteText(filename, text, 0600)
		is.NoError(err1)
		s1, err1 := ReadText(filename)
		is.NoError(err1)
		return s1, nil
	})
	is.NoError(err)
	is.NoDirExists(tmpDir)
	is.Equal(text, s)
}
