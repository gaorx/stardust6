package sdbytes

import (
	"github.com/gaorx/stardust6/sdrand"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncoding(t *testing.T) {
	is := assert.New(t)
	a := sdrand.Bytes(8192)
	is.Equal(a, lo.Must(FromHex(ToHexL(a))))
	is.Equal(a, lo.Must(FromBase64Std(ToBase64Std(a))))
	is.Equal(a, lo.Must(FromBase64Url(ToBase64Url(a))))
}
