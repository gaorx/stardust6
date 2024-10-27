package sderr

import (
	"github.com/stretchr/testify/assert"
	"io/fs"
	"testing"
)

func TestCode(t *testing.T) {
	is := assert.New(t)
	err1 := WithCode("AA").Wrap(fs.ErrExist)
	is.Equal("AA", Code(err1))
	err2 := WithCode("BB").Wrap(err1)
	is.Equal("BB", Code(err2))
}

func TestPublicMsg(t *testing.T) {
	is := assert.New(t)
	err1 := WithPublicMsg("ERROR1").Wrap(fs.ErrExist)
	is.Equal("ERROR1", PublicMsg(err1))
	err2 := WithPublicMsg("ERROR2").Wrap(err1)
	is.Equal("ERROR2", PublicMsg(err2))
	err3 := Wrap(Wrapf(fs.ErrExist, "NO PUBLIC MSG"))
	is.Equal("NO PUBLIC MSG: "+fs.ErrExist.Error(), PublicMsg(err3))
}
