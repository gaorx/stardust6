package sdcrypto

import (
	"bytes"
	"github.com/gaorx/stardust6/sdrand"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAES(t *testing.T) {
	is := assert.New(t)
	data0 := []byte(sdrand.String(1303, sdrand.LowerCaseLettersCharset))
	key := []byte(sdrand.String(16, sdrand.AlphanumericCharset))

	data1, err := AES.Encrypt(key, data0)
	is.NoError(err)
	data2, err := AES.Decrypt(key, data1)
	is.NoError(err)
	is.True(bytes.Equal(data0, data2))

	data3, err := AESCRC32.Encrypt(key, data0)
	is.NoError(err)
	data4, err := AESCRC32.Decrypt(key, data3)
	is.NoError(err)
	is.True(bytes.Equal(data0, data4))
}
