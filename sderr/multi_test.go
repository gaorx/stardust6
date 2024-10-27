package sderr

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/fs"
	"testing"
)

func TestJoin(t *testing.T) {
	is := assert.New(t)

	err1 := fmt.Errorf("ERR1")
	err2 := Newf("%w", fs.ErrClosed)
	err3 := Wrapf(fs.ErrExist, "ERR3")
	all := Join(nil, err1, nil, err2, nil, err3)
	_, isMulti := ProbeMulti(all)
	is.True(isMulti)
	is.True(Is(all, err1))
	is.True(Is(all, err2) && Is(all, fs.ErrClosed))
	is.True(Is(all, fs.ErrExist))
}

func TestMultiError_Error_Lines(t *testing.T) {
	is := assert.New(t)

	err1 := fmt.Errorf("ERR1")
	err2 := Newf("%w", fs.ErrClosed)
	err3 := Wrapf(fs.ErrExist, "ERR3")
	all := Join(nil, err1, nil, err2, nil, err3)
	multiErr, isMulti := ProbeMulti(all)
	is.True(isMulti)
	is.Equal("(ERR1 & file already closed & ERR3: file already exists)", all.Error())

	lines := multiErr.Lines("*")
	is.Equal(3, len(lines))
	is.Equal("*ERR1", lines[0])
	is.Equal("*file already closed", lines[1])
	is.Equal("*ERR3: file already exists", lines[2])
}
