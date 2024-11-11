package sdsha256

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test(t *testing.T) {
	is := assert.New(t)
	const h = "2CF24DBA5FB0A30E26E83B2AC5B9E29E1B161E5C1FA7425E73043362938B9824"
	is.Equal(h, Hash([]byte("hello")).HexU())
	is.True(VerifyHex(h, []byte("hello")))
	is.False(VerifyHex(h, []byte("hell")))

	const hh = "9A7DB05711F6F880984962C3F3955EFFC9CD2472BCEBE3DFDE834314892EE55E"
	const k = "6aYoRHT*5TlX$&8g"
	is.Equal(hh, HashHmac([]byte(k), []byte("hello")).HexU())
	is.True(VerifyHmacHex(hh, []byte(k), []byte("hello")))
	is.False(VerifyHmacHex(hh, []byte(k), []byte("hell")))
}
