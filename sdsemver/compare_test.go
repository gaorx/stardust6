package sdsemver

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCompare(t *testing.T) {
	is := assert.New(t)
	is.True(Compare(New(1, 0, 9), New(2, 0, 4)) < 0)
	is.True(Compare(New(0, 2, 9), New(0, 3, 2)) < 0)
	is.True(Compare(New(0, 0, 3), New(0, 0, 4)) < 0)
	is.True(Compare(New(3, 3, 3), New(3, 3, 3)) == 0)
	is.True(Compare(New(2, 0, 4), New(1, 0, 9)) > 0)
	is.True(Compare(New(0, 3, 2), New(0, 2, 9)) > 0)
	is.True(Compare(New(0, 0, 4), New(0, 0, 3)) > 0)
}
