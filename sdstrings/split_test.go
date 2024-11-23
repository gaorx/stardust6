package sdstrings

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitNonempty(t *testing.T) {
	is := assert.New(t)
	is.Empty(SplitNonempty("", ",", false))
	is.Empty(SplitNonempty(",", ",", false))
	is.Empty(SplitNonempty(",,", ",", false))
	is.EqualValues([]string{"a", "b"}, SplitNonempty(",a,,b", ",", false))
	is.EqualValues([]string{"a", "b"}, SplitNonempty(", a , , b", ",", true))
}

func TestSplit2s(t *testing.T) {
	is := assert.New(t)
	var s1, s2 string

	s1, s2 = Split2s("a.b", ".")
	is.Equal("a", s1)
	is.Equal("b", s2)

	s1, s2 = Split2s("", ".")
	is.Equal("", s1)
	is.Equal("", s2)

	s1, s2 = Split2s("a", ".")
	is.Equal("a", s1)
	is.Equal("", s2)

	s1, s2 = Split2s("a.b.c", ".")
	is.Equal("a", s1)
	is.Equal("b.c", s2)
}

func TestSplit3s(t *testing.T) {
	is := assert.New(t)
	var s1, s2, s3 string

	s1, s2, s3 = Split3s("a.b.c", ".")
	is.Equal("a", s1)
	is.Equal("b", s2)
	is.Equal("c", s3)

	s1, s2, s3 = Split3s("", ".")
	is.Equal("", s1)
	is.Equal("", s2)
	is.Equal("", s3)

	s1, s2, s3 = Split3s("a.b", ".")
	is.Equal("a", s1)
	is.Equal("b", s2)
	is.Equal("", s3)

	s1, s2, s3 = Split3s("a.b.c.d", ".")
	is.Equal("a", s1)
	is.Equal("b", s2)
	is.Equal("c.d", s3)
}

func TestSplit4s(t *testing.T) {
	is := assert.New(t)
	var s1, s2, s3, s4 string

	s1, s2, s3, s4 = Split4s("a.b.c.d", ".")
	is.Equal("a", s1)
	is.Equal("b", s2)
	is.Equal("c", s3)
	is.Equal("d", s4)

	s1, s2, s3, s4 = Split4s("", ".")
	is.Equal("", s1)
	is.Equal("", s2)
	is.Equal("", s3)
	is.Equal("", s4)

	s1, s2, s3, s4 = Split4s("a.b", ".")
	is.Equal("a", s1)
	is.Equal("b", s2)
	is.Equal("", s3)
	is.Equal("", s4)

	s1, s2, s3, s4 = Split4s("a.b.c", ".")
	is.Equal("a", s1)
	is.Equal("b", s2)
	is.Equal("c", s3)
	is.Equal("", s4)

	s1, s2, s3, s4 = Split4s("a.b.c.d.e", ".")
	is.Equal("a", s1)
	is.Equal("b", s2)
	is.Equal("c", s3)
	is.Equal("d.e", s4)
}
