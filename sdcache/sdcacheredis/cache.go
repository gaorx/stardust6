package sdcacheredis

import (
	"context"
	"github.com/gaorx/stardust6/sdcache"
	"github.com/gaorx/stardust6/sderr"
	"github.com/gaorx/stardust6/sdredis"
	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
	"time"
)

// Cache 使用Redis作为后端的缓存实现
type Cache[K sdcache.Key, V any] struct {
	client   redis.UniversalClient
	keyCodec sdcache.Codec[K, string]
	valCodec sdcache.Codec[V, []byte]
	options  *Options
}

// Options 创建缓存时的选项
type Options struct {
	// KeyPrefix 缓存key的前缀
	KeyPrefix string
	// TTL 缓存的默认TTL
	TTL time.Duration
	// OnRedisSet 设置缓存时的回调(成功设置后才执行)
	OnRedisSet func(ctx context.Context, client redis.UniversalClient, key string, val []byte)
	// OnRedisDelete 删除缓存时的回调(成功删除后才执行)
	OnRedisDelete func(ctx context.Context, client redis.UniversalClient, key string)
	// OnRedisClear 清空缓存时的回调(成功晴空后才执行)
	OnRedisClear func(ctx context.Context, client redis.UniversalClient, prefix string)
}

var _ sdcache.Cache[string, string] = Cache[string, string]{}

// New 从一个redis client创建一个Redis缓存
func New[K sdcache.Key, V any](
	client redis.UniversalClient,
	keyCodec sdcache.Codec[K, string],
	valCodec sdcache.Codec[V, []byte],
	opts *Options,
) Cache[K, V] {
	opts1 := lo.FromPtr(opts)
	if opts1.TTL < 0 {
		opts1.TTL = 0
	}
	return Cache[K, V]{
		client:   client,
		keyCodec: keyCodec,
		valCodec: valCodec,
		options:  &opts1,
	}
}

// Dial 从一个地址创建一个Redis缓存
func Dial[K sdcache.Key, V any](
	addr sdredis.Address,
	keyCodec sdcache.Codec[K, string],
	valCodec sdcache.Codec[V, []byte],
	opts *Options,
) (Cache[K, V], error) {
	client, err := sdredis.Dial(addr)
	if err != nil {
		return Cache[K, V]{}, sderr.Wrap(err)
	}
	return New[K, V](client, keyCodec, valCodec, opts), nil
}

// IsZero 是此缓存是否为空
func (c Cache[K, V]) IsZero() bool {
	var zero Cache[K, V]
	return c == zero
}

// Client 获取Redis客户端
func (c Cache[K, V]) Client() redis.UniversalClient {
	return c.client
}

// Options 获取选项
func (c Cache[K, V]) Options() *Options {
	opts := *c.options
	return &opts
}

// Clear 实现 sdcache.Cache 接口
func (c Cache[K, V]) Clear(ctx context.Context) error {
	prefix := c.options.KeyPrefix
	if prefix != "" {
		var cursor uint64 = 0
		for {
			keys, nextCursor, err := c.client.Scan(ctx, cursor, prefix+"*", 1).Result()
			if err != nil {
				return sderr.Wrapf(err, "scan redis key error")
			}
			if len(keys) <= 0 {
				break
			}
			err = c.client.Del(ctx, keys...).Err()
			if err != nil {
				return sderr.Wrapf(err, "delete redis key error")
			}
			cursor = nextCursor
		}
	} else {
		err := c.client.FlushAll(ctx).Err()
		if err != nil {
			return sderr.Wrapf(err, "clear redis data error")
		}
	}
	c.fireOnRedisClear(ctx, c.client, prefix)
	return nil
}

// Get 实现 sdcache.Cache 接口
func (c Cache[K, V]) Get(ctx context.Context, key K) (V, error) {
	var zero V
	encodedKey, err := c.getKey(key)
	if err != nil {
		return zero, sderr.Wrapf(err, "encode key failed")
	}
	valueBytes, err := c.client.Get(ctx, encodedKey).Bytes()
	if err != nil {
		if sderr.Is(err, redis.Nil) {
			return zero, sderr.Wrapf(sdcache.ErrNotFound, "get redis key error")
		} else {
			return zero, sderr.Wrapf(err, "get redis key error")
		}
	}
	r, err := c.valCodec.Decode(valueBytes)
	if err != nil {
		return zero, sderr.Wrapf(err, "decode value failed")
	}
	return r, nil
}

// GetTTL 实现 sdcache.Cache 接口
func (c Cache[K, V]) GetTTL(ctx context.Context, key K) (time.Duration, error) {
	encodedKey, err := c.getKey(key)
	if err != nil {
		return 0, sderr.Wrapf(err, "encode key failed")
	}
	ttl, err := c.client.TTL(ctx, encodedKey).Result()
	if err != nil {
		return 0, sderr.Wrapf(err, "get ttl failed")
	}
	if ttl == -1 {
		return 0, nil
	} else if ttl == -2 {
		return 0, sderr.Wrapf(sdcache.ErrNotFound, "get redis key ttl error")
	} else {
		return ttl, nil
	}
}

// Put 实现 sdcache.Cache 接口
func (c Cache[K, V]) Put(ctx context.Context, key K, val V, opts *sdcache.PutOptions) error {
	encodedKey, err := c.getKey(key)
	if err != nil {
		return sderr.Wrapf(err, "encode key failed")
	}
	encodedValue, err := c.valCodec.Encode(val)
	if err != nil {
		return sderr.Wrapf(err, "encode value failed")
	}
	err = c.client.Set(ctx, encodedKey, encodedValue, c.getTTL(opts)).Err()
	if err != nil {
		return sderr.Wrapf(err, "set redis value error")
	}
	c.fireOnRedisSet(ctx, c.client, encodedKey, encodedValue)
	return nil
}

// Delete 实现 sdcache.Cache 接口
func (c Cache[K, V]) Delete(ctx context.Context, key K) error {
	encodedKey, err := c.getKey(key)
	if err != nil {
		return sderr.Wrapf(err, "encode key failed")
	}
	err = c.client.Del(ctx, encodedKey).Err()
	if err != nil {
		return sderr.Wrapf(err, "delete redis key error")
	}
	c.fireOnRedisDelete(ctx, c.client, encodedKey)
	return nil
}

// GetOrLoad 实现 sdcache.Cache 接口
func (c Cache[K, V]) GetOrLoad(ctx context.Context, key K, loader sdcache.Loader[K, V], opts *sdcache.PutOptions) (V, error) {
	var zero V
	if loader == nil {
		return zero, sderr.Newf("nil loader")
	}
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
	encodedKey, err := c.getKey(key)
	if err != nil {
		return zero, sderr.Wrapf(err, "encode key failed")
	}
	encodedValue, err := c.client.Get(ctx, encodedKey).Bytes()
	if err != nil {
		if sderr.Is(err, redis.Nil) {
			v, found, err := loader(ctx, key)
			if err != nil {
				return zero, sderr.Wrapf(err, "load value for redis error")
			}
			if !found {
				return zero, sderr.Wrapf(sdcache.ErrNotFound, "load nothing")
			}
			encodedValue1, err := c.valCodec.Encode(v)
			if err != nil {
				return zero, sderr.Wrapf(err, "encode value failed")
			}
			err = c.client.Set(ctx, encodedKey, encodedValue1, c.getTTL(opts)).Err()
			if err != nil {
				return zero, sderr.Wrapf(err, "set redis value error")
			}
			c.fireOnRedisSet(ctx, c.client, encodedKey, encodedValue1)
			return v, nil
		} else {
			return zero, sderr.Wrapf(err, "get value in redis failed")
		}
	} else {
		v, err := c.valCodec.Decode(encodedValue)
		if err != nil {
			return zero, sderr.Wrapf(err, "decode redis value error")
		}
		return v, nil
	}
}

func (c Cache[K, V]) getKey(key K) (string, error) {
	encodedKey, err := c.keyCodec.Encode(key)
	if err != nil {
		return "", err
	}
	return c.options.KeyPrefix + encodedKey, nil
}

func (c Cache[K, V]) getTTL(opts *sdcache.PutOptions) time.Duration {
	if opts != nil && opts.TTL > 0 {
		return opts.TTL
	}
	return c.options.TTL
}

func (c Cache[K, V]) fireOnRedisSet(ctx context.Context, client redis.UniversalClient, key string, val []byte) {
	f := c.options.OnRedisSet
	if f == nil {
		return
	}
	f(ctx, client, key, val)
}

func (c Cache[K, V]) fireOnRedisDelete(ctx context.Context, client redis.UniversalClient, key string) {
	f := c.options.OnRedisDelete
	if f == nil {
		return
	}
	f(ctx, client, key)
}

func (c Cache[K, V]) fireOnRedisClear(ctx context.Context, client redis.UniversalClient, prefix string) {
	f := c.options.OnRedisClear
	if f == nil {
		return
	}
	f(ctx, client, prefix)
}
