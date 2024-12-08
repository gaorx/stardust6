package sdcacheredis

import (
	"github.com/gaorx/stardust6/sdcache"
	"github.com/gaorx/stardust6/sdcache/internal"
	"github.com/gaorx/stardust6/sdredis"
	"github.com/gaorx/stardust6/sdtime"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCache(t *testing.T) {
	is := assert.New(t)
	ttlSecs := int64(2)
	c, err := Dial[string, string](
		sdredis.Address{
			Addrs:    []string{"localhost:6379"},
			Password: "",
			DB:       1,
		},
		sdcache.StringToString(),
		sdcache.StringToBytes(),
		&Options{
			TTL: sdtime.Seconds(ttlSecs),
		},
	)
	is.NoError(err)

	// tests
	internal.TestCommon(t, c)
	internal.TestExpiration(t, c, ttlSecs)
}
