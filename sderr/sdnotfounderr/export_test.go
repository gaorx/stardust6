package sdnotfounderr

import (
	"database/sql"
	"github.com/gaorx/stardust6/sdcache"
	"github.com/gaorx/stardust6/sderr"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIs(t *testing.T) {
	is := assert.New(t)

	// 自定义not found 语义error
	errCustom := sderr.Newf("custom not found")

	is.True(Is(sdcache.ErrNotFound))
	is.True(Is(redis.Nil))
	is.True(Is(sql.ErrNoRows))
	is.False(Is(errCustom))

	Register(errCustom)
	is.True(Is(sdcache.ErrNotFound))
	is.True(Is(redis.Nil))
	is.True(Is(sql.ErrNoRows))
	is.True(Is(errCustom))
}
