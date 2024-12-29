package sdcodegen

import (
	"github.com/gaorx/stardust6/sdfile"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func TestDir(t *testing.T) {
	is := assert.New(t)
	_ = sdfile.UseTempDir("", "", func(dirname string) {
		projDir := filepath.Join(dirname, "project1")
		g0 := New()
		g0.Add("README.md", Text("TODO"))
		g0.Add(".gitignore", Text("*.log"))
		g0.Add("src/a.js", Text("console.log('a')"))
		err := g0.Generate(projDir)
		is.NoError(err)

		g1 := New()
		g1.AddModule(Dir(projDir))
		g1.Add("other.txt", Text("other"))
		outs, err := g1.TryAll(nil)
		is.NoError(err)
		is.Len(outs, 4)
		is.Equal("TODO", outs.Get("README.md").Text())
		is.Equal("*.log", outs.Get(".gitignore").Text())
		is.Equal("console.log('a')", outs.Get("src/a.js").Text())
		is.Equal("other", outs.Get("other.txt").Text())
	})
}
