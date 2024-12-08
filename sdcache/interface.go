package sdcache

import (
	"context"
	"time"
)

// Cache 缓存接口
type Cache[K Key, V any] interface {
	// Clear 清空缓存
	Clear(ctx context.Context) error
	// Get 获取某个缓存key的所对应的值，错误的话返回error，如果不存在返回err的cause是ErrNotFound
	Get(ctx context.Context, key K) (V, error)
	// GetTTL 获取某个缓存key的TTL，如果不存在返回error
	GetTTL(ctx context.Context, key K) (time.Duration, error)
	// Put 将某个值放入缓存，如果已经存在则覆盖，如果opts中的TTL为空，则使用默认值
	Put(ctx context.Context, key K, val V, opts *PutOptions) error
	// Delete 删除某个缓存key
	Delete(ctx context.Context, key K) error
	// GetOrLoad 获取某个缓存key的值，如果不存在则调用loader加载，并将加载成功的结果放入缓存中，
	// 如果opts中的TTL为空，则使用默认值
	GetOrLoad(ctx context.Context, key K, loader Loader[K, V], opts *PutOptions) (V, error)
}

// Key 缓存key接口
type Key interface {
	string | []byte | int64 | uint64
}

// PutOptions 缓存操作选项
type PutOptions struct {
	// TTL 此次Put的TTL，如果为0则使用默认值
	TTL time.Duration
	// Cost 此次Put的成本，如果为0则使用默认值
	Cost int64
}

// Loader 缓存值加载器
// 如果加载错误，返回error；如果不存在，返回false；否则返回true
type Loader[K Key, V any] func(ctx context.Context, key K) (V, bool, error)
