package sdmd5

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test(t *testing.T) {
	is := assert.New(t)
	const h = "5D41402ABC4B2A76B9719D911017C592"
	is.Equal(h, Hash([]byte("hello")).HexU())
	is.True(VerifyHex(h, []byte("hello")))
	is.False(VerifyHex(h, []byte("hell")))
}
