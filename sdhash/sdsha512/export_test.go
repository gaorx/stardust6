package sdsha512

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test(t *testing.T) {
	is := assert.New(t)
	const h = "9B71D224BD62F3785D96D46AD3EA3D73319BFBC2890CAADAE2DFF72519673CA72323C3D99BA5C11D7C7ACC6E14B8C5DA0C4663475C2E5C3ADEF46F73BCDEC043"
	is.Equal(h, Hash([]byte("hello")).HexU())
	is.True(VerifyHex(h, []byte("hello")))
	is.False(VerifyHex(h, []byte("hell")))

	const hh = "427A3F3F94D7598F2FFE5767D3D51B419FEAE64562A52A18FFB0A93CCDFF67D26D667B0BD93E303D943824D882D0185A98FC255D872C82BCABB4279331452B52"
	const k = "6aYoRHT*5TlX$&8g"
	is.Equal(hh, HashHmac([]byte(k), []byte("hello")).HexU())
	is.True(VerifyHmacHex(hh, []byte(k), []byte("hello")))
	is.False(VerifyHmacHex(hh, []byte(k), []byte("hell")))
}
