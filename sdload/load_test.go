package sdload

import (
	"github.com/gaorx/stardust6/sdfile"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoad(t *testing.T) {
	is := assert.New(t)

	_ = sdfile.UseTempDir("", "", func(dirname string) {
		const text = "hello world"
		fn := filepath.Join(dirname, "a.txt")
		err := sdfile.WriteText(fn, text, 0600)
		is.NoError(err)
		s, err := Text(fn)
		is.NoError(err)
		is.Equal(text, s)
		s, err = Text("file://" + fn)
		is.NoError(err)
		is.Equal(text, s)
	})
	s, err := Text("https://www.baidu.com")
	is.NoError(err)
	is.True(strings.Contains(s, "baidu"))
	s, err = Text("http://www.baidu.com")
	is.NoError(err)
	is.True(strings.Contains(s, "baidu"))
	_, err = Text("https://unavaiable_hostname/aaa.txt")
	is.Error(err)
}
