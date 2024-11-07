package sdcodegen

import (
	"bytes"
	"github.com/gaorx/stardust6/sdfile"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestFIle(t *testing.T) {
	is := assert.New(t)

	is.True((&File{}).IsZero() && (*File)(nil).IsZero() && !(&File{Name: "a.txt"}).IsZero())

	_ = sdfile.UseTempDir("", "", func(dirname string) {
		f1 := &File{
			Name:      "scripts/a.sh",
			Data:      []byte("#!/bin/bash\necho 'hello'\n"),
			Mode:      os.FileMode(0755),
			Discarded: true,
		}
		err := f1.Write(dirname, nil)
		is.NoError(err)
		is.NoFileExists(filepath.Join(dirname, "scripts/a.sh"))
		f1.Discarded = false
		err = f1.Write(dirname, nil)
		is.NoError(err)
		is.FileExists(filepath.Join(dirname, "scripts/a.sh"))

		f2, err := ReadFile(dirname, "scripts/a.sh")
		is.NoError(err)
		is.True(f1.Name == f2.Name && f1.Mode == f2.Mode && bytes.Equal(f1.Data, f2.Data))

		f2, err = ReadFileFS(os.DirFS(dirname), "scripts/a.sh")
		is.NoError(err)
		is.True(f1.Name == f2.Name && f1.Mode == f2.Mode && bytes.Equal(f1.Data, f2.Data))

		_, err = ReadFile(dirname, "not-exists.txt")
		is.Error(err)
		is.True(IsNotExistErr(err))
	})
}
