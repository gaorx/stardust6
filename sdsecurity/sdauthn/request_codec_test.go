package sdauthn

import (
	"github.com/gaorx/stardust6/sdrand"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRequestCodecFunc(t *testing.T) {
	is := assert.New(t)

	codec := JWTUserToken(sdrand.String(16, sdrand.LowerCaseAlphanumericCharset))
	is.Equal("", codec.Encode(NewPrincipalId(PrincipalUid, "domain", "uid")))

	token1 := NewUserToken("uid", time.Now()).In("domain")
	s1, s2 := codec.Encode(token1), codec.Encode(&token1)
	is.True(s1 == s2)
	token2, err := codec.Decode(s1)
	is.NoError(err)
	is.IsType(UserToken{}, token2)
	is.True(token1 == token2)
}
