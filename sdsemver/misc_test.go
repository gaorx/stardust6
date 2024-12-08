package sdsemver

import (
	"github.com/gaorx/stardust6/sdrand"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConvert(t *testing.T) {
	is := assert.New(t)
	for i := 0; i < 10000; i++ {
		major := sdrand.IntBetween(0, numLimit)
		minor := sdrand.IntBetween(0, numLimit)
		patch := sdrand.IntBetween(0, numLimit)
		s0 := New(major, minor, patch).String()
		vi, err := ToInt(s0)
		is.NoError(err)
		s1, err := ToString(vi)
		is.NoError(err)
		is.True(s0 == s1)
	}
}
