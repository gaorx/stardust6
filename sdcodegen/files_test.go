package sdcodegen

import (
	"github.com/gaorx/stardust6/sdfile"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestFiles(t *testing.T) {
	is := assert.New(t)

	_ = sdfile.UseTempDir("", "", func(dirname string) {
		var files Files
		f, err := Generate("src/a.txt", nil, Text("hello\n"))
		is.NoError(err)
		files = append(files, f)

		f, err = Generate("src/b.md", nil, Text("# Title\n"))
		is.NoError(err)
		files = append(files, f)

		err = files.Write(dirname, nil)
		is.NoError(err)

		files, err = ReadDir(dirname)
		is.NoError(err)
		is.Equal(2, len(files))
		is.NotNil(files.Get("src/a.txt"))
		is.NotNil(files.Get("src/b.md"))

		a := files.Get("src/a.txt")
		is.True(a.Name == "src/a.txt")
		is.Equal([]byte("hello\n"), a.Data)
		is.True(a.Mode == 0600)
		is.NotZero(a.ModTime)

		files, err = ReadDirFS(os.DirFS(dirname), HasSuffix(".txt"))
		is.NoError(err)
		is.Equal(1, len(files))
		is.NotNil(files.Get("src/a.txt"))
		is.Nil(files.Get("src/b.md"))
	})
}
