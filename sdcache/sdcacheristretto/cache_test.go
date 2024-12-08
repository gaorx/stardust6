package sdcacheristretto

import (
	"github.com/dgraph-io/ristretto/v2"
	"github.com/gaorx/stardust6/sdcache/internal"
	"github.com/gaorx/stardust6/sdtime"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCache(t *testing.T) {
	is := assert.New(t)
	ttlSecs := int64(2)
	c, err := NewByConfig[string, string](
		&ristretto.Config[string, string]{
			NumCounters: 100,
			MaxCost:     100,
			BufferItems: 64,
		},
		&Options{
			TTL: sdtime.Seconds(ttlSecs),
		},
	)
	is.NoError(err)

	// tests
	internal.TestCommon(t, c)
	internal.TestExpiration(t, c, ttlSecs)
}
