package sdcacheristretto

import (
	"context"
	"github.com/dgraph-io/ristretto/v2"
	"github.com/gaorx/stardust6/sdcache"
	"github.com/gaorx/stardust6/sderr"
	"github.com/samber/lo"
	"time"
)

// Cache 使用Ristretto作为后端的缓存实现
type Cache[K sdcache.Key, V any] struct {
	cache   *ristretto.Cache[K, V]
	options *Options
}

// Options 创建缓存时的选项
type Options struct {
	// 默认的TTL
	TTL time.Duration
	// 默认的Cost
	Cost int64
}

// New 从一个ristretto cache创建一个Ristretto缓存
func New[K sdcache.Key, V any](c *ristretto.Cache[K, V], opts *Options) Cache[K, V] {
	opts1 := lo.FromPtr(opts)
	return Cache[K, V]{
		cache:   c,
		options: &opts1,
	}
}

// NewByConfig 从一个ristretto config创建一个Ristretto缓存
func NewByConfig[K sdcache.Key, V any](config *ristretto.Config[K, V], opts *Options) (Cache[K, V], error) {
	c, err := ristretto.NewCache[K, V](config)
	if err != nil {
		return Cache[K, V]{}, err
	}
	return New[K, V](c, opts), nil
}

// IsZero 判断此缓存是否为零值
func (c Cache[K, V]) IsZero() bool {
	var zero Cache[K, V]
	return c == zero
}

// Ristretto 获取底层的ristretto.Cache
func (c Cache[K, V]) Ristretto() *ristretto.Cache[K, V] {
	return c.cache
}

// Options 获取选项
func (c Cache[K, V]) Options() *Options {
	opts := *c.options
	return &opts
}

// Clear 实现 sdcache.Cache 接口
func (c Cache[K, V]) Clear(_ context.Context) error {
	c.cache.Clear()
	return nil
}

// Get 实现 sdcache.Cache 接口
func (c Cache[K, V]) Get(_ context.Context, key K) (V, error) {
	v, found := c.cache.Get(key)
	if !found {
		var zero V
		return zero, sderr.Wrapf(sdcache.ErrNotFound, "get cache key error")
	}
	return v, nil
}

// GetTTL 实现 sdcache.Cache 接口
func (c Cache[K, V]) GetTTL(_ context.Context, key K) (time.Duration, error) {
	ttl, found := c.cache.GetTTL(key)
	if !found {
		return 0, sderr.Wrapf(sdcache.ErrNotFound, "get cache key ttl error")
	}
	return ttl, nil
}

// Put 实现 sdcache.Cache 接口
func (c Cache[K, V]) Put(_ context.Context, key K, val V, opts *sdcache.PutOptions) error {
	ttl, cost := c.getTTL(opts), c.getCost(opts)
	if ttl > 0 {
		_ = c.cache.SetWithTTL(key, val, cost, ttl)
	} else {
		_ = c.cache.Set(key, val, cost)
	}
	c.cache.Wait() // ristretto v2和v1不同，写入后必须调用这个才能让下一次读直接生效
	return nil
}

// Delete 实现 sdcache.Cache 接口
func (c Cache[K, V]) Delete(_ context.Context, key K) error {
	c.cache.Del(key)
	return nil
}

// GetOrLoad 实现 sdcache.Cache 接口
func (c Cache[K, V]) GetOrLoad(ctx context.Context, key K, loader sdcache.Loader[K, V], opts *sdcache.PutOptions) (V, error) {
	var zero V
	if c.IsZero() {
		v, found, err := loader(ctx, key)
		if err != nil {
			return zero, sderr.Wrapf(err, "load value error")
		}
		if !found {
			return zero, sderr.Wrapf(sdcache.ErrNotFound, "load nothing")
		}
		return v, nil
	}
	v, found := c.cache.Get(key)
	if !found {
		v0, found0, err := loader(ctx, key)
		if err != nil {
			return zero, sderr.Wrapf(err, "load value error")
		}
		if !found0 {
			return zero, sderr.Wrapf(sdcache.ErrNotFound, "load nothing")
		}
		ttl, cost := c.getTTL(opts), c.getCost(opts)
		if ttl > 0 {
			_ = c.cache.SetWithTTL(key, v0, cost, ttl)
		} else {
			_ = c.cache.Set(key, v0, cost)
		}
		c.cache.Wait() // ristretto v2和v1不同，写入后必须调用这个才能让下一次读直接生效
		return v0, nil
	} else {
		return v, nil
	}
}

func (c Cache[K, V]) getTTL(opts *sdcache.PutOptions) time.Duration {
	if opts != nil && opts.TTL > 0 {
		return opts.TTL
	}
	return c.options.TTL
}

func (c Cache[K, V]) getCost(opts *sdcache.PutOptions) int64 {
	if opts == nil || opts.Cost < 0 {
		return c.options.Cost
	}
	return opts.Cost
}
