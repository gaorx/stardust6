package sdauthn

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestExpiration(t *testing.T) {
	is := assert.New(t)
	expiration := ExpireAt(time.Now().Add(200 * time.Millisecond))
	is.False(expiration.IsExpired(time.Now(), time.Now().Add(40*time.Millisecond)))
	is.True(expiration.IsExpired(time.Now(), time.Now().Add(300*time.Millisecond)))

	expiration = ExpireIn(200 * time.Millisecond)
	is.False(expiration.IsExpired(time.Now(), time.Now().Add(40*time.Millisecond)))
	is.True(expiration.IsExpired(time.Now(), time.Now().Add(300*time.Millisecond)))
}
