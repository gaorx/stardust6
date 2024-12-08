package sdsql

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPage(t *testing.T) {
	is := assert.New(t)

	// Page1
	p := Page1(2, 4)
	is.Equal(int64(2), p.Num())
	is.Equal(int64(4), p.Size())
	limit, offset := p.LimitAndOffset()
	is.Equal(int64(4), limit)
	is.Equal(int64(4), offset)

	p = Page1(-1, 0)
	is.Equal(int64(1), p.Num())
	is.Equal(int64(1), p.Size())
	limit, offset = p.LimitAndOffset()
	is.Equal(int64(1), limit)
	is.Equal(int64(0), offset)

	p = Page1(1, 1000000000)
	is.Equal(int64(1), p.Num())
	is.Equal(int64(100000), p.Size())
	limit, offset = p.LimitAndOffset()
	is.Equal(int64(100000), limit)
	is.Equal(int64(0), offset)

	// Page0
	p = Page0(2, 4)
	is.Equal(int64(2), p.Num())
	is.Equal(int64(4), p.Size())
	limit, offset = p.LimitAndOffset()
	is.Equal(int64(4), limit)
	is.Equal(int64(8), offset)

	p = Page0(-1, 0)
	is.Equal(int64(0), p.Num())
	is.Equal(int64(1), p.Size())
	limit, offset = p.LimitAndOffset()
	is.Equal(int64(1), limit)
	is.Equal(int64(0), offset)

	p = Page0(1, 1000000000)
	is.Equal(int64(1), p.Num())
	is.Equal(int64(100000), p.Size())
	limit, offset = p.LimitAndOffset()
	is.Equal(int64(100000), limit)
	is.Equal(int64(100000), offset)
}
