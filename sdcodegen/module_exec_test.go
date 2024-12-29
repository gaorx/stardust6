package sdcodegen

import (
	"fmt"
	"github.com/gaorx/stardust6/sdexec"
	"github.com/gaorx/stardust6/sdfile"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestExecForDir(t *testing.T) {
	is := assert.New(t)
	_ = sdfile.UseTempDir("", "", func(dirname string) {
		d := filepath.Join(dirname, "project1")
		err := os.MkdirAll(d, 0755)
		is.NoError(err)
		line := fmt.Sprintf(`echo -n "hello" > "%s/hello.txt"`, d)
		g1 := New()
		g1.Sub("sub").AddModule(ExecForDir(sdexec.Bash(line), nil, d))
		outs, err := g1.TryAll(nil)
		is.NoError(err)
		is.Len(outs, 1)
		is.Equal("hello", outs.Get("sub/hello.txt").Text())
	})
}
