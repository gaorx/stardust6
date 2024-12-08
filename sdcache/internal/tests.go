package internal

import (
	"context"
	"github.com/gaorx/stardust6/sdcache"
	"github.com/gaorx/stardust6/sderr"
	"github.com/gaorx/stardust6/sdtime"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCommon(t *testing.T, c sdcache.Cache[string, string]) {
	is := assert.New(t)

	// clear
	err := c.Clear(context.Background())
	is.NoError(err)

	// Get
	v1, err := c.Get(context.Background(), "k1")
	is.ErrorIs(err, sdcache.ErrNotFound)
	is.Zero(v1)

	// Put
	err = c.Put(context.Background(), "k1", "v1", nil)
	is.NoError(err)

	// Get
	v1, err = c.Get(context.Background(), "k1")
	is.NoError(err)
	is.Equal("v1", v1)

	// Delete
	err = c.Delete(context.Background(), "k1")
	is.NoError(err)

	// Get
	v1, err = c.Get(context.Background(), "k1")
	is.ErrorIs(err, sdcache.ErrNotFound)
	is.Zero(v1)

	// GetOrLoad
	var errLoad = sderr.Newf("load error")
	loadCounter := 0
	loader := func(ctx context.Context, k string) (string, bool, error) {
		loadCounter += 1
		switch k {
		case "k1":
			return "v1", true, nil // 命中
		case "k2":
			return "", false, errLoad // 模拟加载失败
		default:
			return "", false, nil // 未命中，但并没有加载失败
		}
	}

	// k2
	v2, err := c.GetOrLoad(context.Background(), "k2", loader, nil)
	is.ErrorIs(err, errLoad)
	is.Zero(v2)
	is.Equal(1, loadCounter)
	v2, err = c.GetOrLoad(context.Background(), "k2", loader, nil)
	is.ErrorIs(err, errLoad)
	is.Zero(v2)
	is.Equal(2, loadCounter)
	v2, err = c.Get(context.Background(), "k2")
	is.ErrorIs(err, sdcache.ErrNotFound)
	is.Zero(v2)

	// k3
	v3, err := c.GetOrLoad(context.Background(), "k3", loader, nil)
	is.ErrorIs(err, sdcache.ErrNotFound)
	is.Zero(v3)
	is.Equal(3, loadCounter)
	v3, err = c.GetOrLoad(context.Background(), "k3", loader, nil)
	is.ErrorIs(err, sdcache.ErrNotFound)
	is.Zero(v3)
	is.Equal(4, loadCounter)
	v3, err = c.Get(context.Background(), "k3")
	is.ErrorIs(err, sdcache.ErrNotFound)
	is.Zero(v3)

	// k1
	v1, err = c.GetOrLoad(context.Background(), "k1", loader, nil)
	is.NoError(err)
	is.Equal("v1", v1)
	is.Equal(5, loadCounter)
	v1, err = c.GetOrLoad(context.Background(), "k1", loader, nil)
	is.NoError(err)
	is.Equal("v1", v1)
	is.Equal(5, loadCounter)
	v1, err = c.Get(context.Background(), "k1")
	is.NoError(err)
	is.Equal("v1", v1)
}

func TestExpiration(t *testing.T, c sdcache.Cache[string, string], ttlSecs int64) {
	is := assert.New(t)
	is.True(ttlSecs > 1)

	// clear
	err := c.Clear(context.Background())
	is.NoError(err)

	// Get
	v1, err := c.Get(context.Background(), "k1")
	is.ErrorIs(err, sdcache.ErrNotFound)
	is.Zero(v1)

	// Put
	err = c.Put(context.Background(), "k1", "v1", &sdcache.PutOptions{
		TTL: sdtime.Seconds(ttlSecs + 1),
	})
	is.NoError(err)
	sdtime.SleepS(1)
	v1, err = c.Get(context.Background(), "k1")
	is.NoError(err)
	is.Equal("v1", v1)
	sdtime.SleepS(ttlSecs + 1)
	_, err = c.Get(context.Background(), "k1")
	is.ErrorIs(err, sdcache.ErrNotFound)

	// GetOrLoad
	var errLoad = sderr.Newf("load error")
	loadCounter := 0
	loader := func(ctx context.Context, k string) (string, bool, error) {
		loadCounter += 1
		switch k {
		case "k1":
			return "v1", true, nil // 命中
		case "k2":
			return "", false, errLoad // 模拟加载失败
		default:
			return "", false, nil // 未命中，但并没有加载失败
		}
	}

	// k2
	v2, err := c.GetOrLoad(context.Background(), "k2", loader, nil)
	is.ErrorIs(err, errLoad)
	is.Zero(v2)
	is.Equal(1, loadCounter)
	v2, err = c.Get(context.Background(), "k2")
	is.ErrorIs(err, sdcache.ErrNotFound)
	is.Zero(v2)

	// k1
	v1, err = c.GetOrLoad(context.Background(), "k1", loader, nil)
	is.NoError(err)
	is.Equal("v1", v1)
	is.Equal(2, loadCounter)
	sdtime.SleepS(1)
	v1, err = c.GetOrLoad(context.Background(), "k1", loader, nil)
	is.NoError(err)
	is.Equal("v1", v1)
	is.Equal(2, loadCounter)
	sdtime.SleepS(ttlSecs)

	// k1
	v1, err = c.Get(context.Background(), "k1")
	is.ErrorIs(err, sdcache.ErrNotFound)
	v1, err = c.GetOrLoad(context.Background(), "k1", loader, &sdcache.PutOptions{
		TTL: sdtime.Seconds(2),
	})
	is.NoError(err)
	is.Equal("v1", v1)
	is.Equal(3, loadCounter)
}
