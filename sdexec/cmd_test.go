package sdexec

import (
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
)

func TestCmd(t *testing.T) {
	is := assert.New(t)

	msg := "hello world!"
	homeDir, err := os.UserHomeDir()
	is.NoError(err)
	cmd, err := Parsef("echo '%s'", msg)
	is.NoError(err)
	cmd.SetDir(homeDir)
	rr := cmd.RunResult()
	is.NoError(rr.Err)
	is.True(strings.Contains(rr.StdoutString(), msg))
}
