package sdsemver

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParse(t *testing.T) {
	is := assert.New(t)
	_, err := Parse("")
	is.Error(err)
	_, err = Parse("a.b.c")
	is.Error(err)
	_, err = Parse("1000000.1.1")
	is.Error(err)
	_, err = Parse("1.10000000.1")
	is.Error(err)
	_, err = Parse("1.0.1000000")
	is.Error(err)
	v, err := Parse("3")
	is.True(v.Equal(3, 0, 0))
	v, err = Parse("0.3")
	is.True(v.Equal(0, 3, 0))
	v, err = Parse("0.2.3")
	is.True(v.Equal(0, 2, 3))
	_, err = Parse("0.2.3.4")
	is.Error(err)
}
