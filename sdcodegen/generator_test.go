package sdcodegen

import (
	"bytes"
	"github.com/gaorx/stardust6/sdfile"
	"github.com/stretchr/testify/assert"
	"io/fs"
	"testing"
)

func TestGenerator1(t *testing.T) {
	is := assert.New(t)

	fileIs := func(f *File, name string, mode fs.FileMode, data []byte, discarded bool) bool {
		if f == nil {
			return false
		}
		return f.Name == name && f.Mode == mode && bytes.Equal(f.Data, data) && f.Discarded == discarded
	}

	g := New()
	g.Add("a.txt", Text("hello"))
	g.Add("b.txt", Text("world"))

	all, err := g.TryAll(nil)
	is.NoError(err)
	is.Equal(2, len(all))
	is.True(fileIs(all.Get("a.txt"), "a.txt", 0600, []byte("hello"), false))
	is.True(fileIs(all.Get("b.txt"), "b.txt", 0600, []byte("world"), false))

	one, err := g.TryOne("b.txt", nil)
	is.NoError(err)
	is.True(fileIs(one, "b.txt", 0600, []byte("world"), false))
	_, err = g.TryOne("not-exists.txt", nil)
	is.Error(err)
}

func TestGenerator2(t *testing.T) {
	is := assert.New(t)

	fileIs := func(f *File, name string, mode fs.FileMode, data []byte) bool {
		if f == nil {
			return false
		}
		return f.Name == name && f.Mode == mode && bytes.Equal(f.Data, data)
	}

	g := New()
	g.Add("a.txt", func(c *Context) {
		c.DiscardAndAbortIfExists()
		if c.Current() != nil {
			c.Line("WORLD")
		} else {
			c.Line("HELLO")
		}
	})
	g.Add("b.txt", func(c *Context) {
		if c.Current() != nil {
			c.Line("WORLD")
		} else {
			c.Line("HELLO")
		}
	})
	g.Add("c.txt", func(c *Context) {
		c.SetExecutable(true)
		c.Line("HELLO")
	})

	_ = sdfile.UseTempDir("", "", func(dirname string) {
		err := g.Generate(dirname)
		is.NoError(err)
		f, err := ReadFile(dirname, "a.txt")
		is.NoError(err)
		is.True(fileIs(f, "a.txt", 0600, []byte("HELLO\n")))
		f, err = ReadFile(dirname, "b.txt")
		is.NoError(err)
		is.True(fileIs(f, "b.txt", 0600, []byte("HELLO\n")))
		f, err = ReadFile(dirname, "c.txt")
		is.NoError(err)
		is.True(fileIs(f, "c.txt", 0700, []byte("HELLO\n")))

		err = g.Generate(dirname)
		is.NoError(err)
		f, err = ReadFile(dirname, "a.txt")
		is.NoError(err)
		is.True(fileIs(f, "a.txt", 0600, []byte("HELLO\n")))
		f, err = ReadFile(dirname, "b.txt")
		is.NoError(err)
		is.True(fileIs(f, "b.txt", 0600, []byte("WORLD\n")))
		f, err = ReadFile(dirname, "c.txt")
		is.NoError(err)
		is.True(fileIs(f, "c.txt", 0700, []byte("HELLO\n")))
	})
}

func TestGenerate3(t *testing.T) {
	is := assert.New(t)

	g := New()
	g.Add("a.txt", func(c *Context) {
		if c.Current() != nil {
			c.Line("WORLD")
		} else {
			c.Line("HELLO")
		}
	})
	_ = sdfile.UseTempDir("", "", func(dirname string) {
		err := g.GenerateOne(dirname, "a.txt")
		is.NoError(err)
		f, err := ReadFile(dirname, "a.txt")
		is.NoError(err)
		is.Equal("HELLO\n", string(f.Data))

		err = g.GenerateOne(dirname, "a.txt")
		is.NoError(err)
		f, err = ReadFile(dirname, "a.txt")
		is.NoError(err)
		is.Equal("WORLD\n", string(f.Data))
	})
}
