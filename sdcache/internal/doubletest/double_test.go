package doubletest

import (
	"github.com/dgraph-io/ristretto/v2"
	"github.com/gaorx/stardust6/sdcache"
	"github.com/gaorx/stardust6/sdcache/internal"
	"github.com/gaorx/stardust6/sdcache/sdcacheredis"
	"github.com/gaorx/stardust6/sdcache/sdcacheristretto"
	"github.com/gaorx/stardust6/sdredis"
	"github.com/gaorx/stardust6/sdtime"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDouble(t *testing.T) {
	is := assert.New(t)
	ttlSecs := int64(2)

	// 一级缓存
	c1, err := sdcacheredis.Dial[string, string](
		sdredis.Address{
			Addrs:    []string{"localhost:6379"},
			Password: "",
			DB:       1,
		},
		sdcache.StringToString(),
		sdcache.StringToBytes(),
		&sdcacheredis.Options{
			TTL: sdtime.Seconds(ttlSecs),
		},
	)
	is.NoError(err)

	// 二级缓存
	c2, err := sdcacheristretto.NewByConfig[string, string](
		&ristretto.Config[string, string]{
			NumCounters: 100,
			MaxCost:     100,
			BufferItems: 64,
		},
		&sdcacheristretto.Options{
			TTL: sdtime.Seconds(ttlSecs),
		},
	)
	is.NoError(err)

	// 双层缓存
	d12 := sdcache.Double[string, string](c1, c2)
	d1 := sdcache.Double[string, string](c1, nil)
	d2 := sdcache.Double[string, string](nil, c2)

	// tests
	internal.TestCommon(t, d12)
	internal.TestExpiration(t, d12, ttlSecs)
	internal.TestCommon(t, d1)
	internal.TestExpiration(t, d1, ttlSecs)
	internal.TestCommon(t, d2)
	internal.TestExpiration(t, d2, ttlSecs)
}
