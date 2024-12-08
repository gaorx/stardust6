package sdredis

import (
	"context"
	"github.com/gaorx/stardust6/sderr"
	"github.com/redis/go-redis/v9"
)

// ForEachShards 遍历client中所有分片，每个分片是一个redis.Client
func ForEachShards(ctx context.Context, client redis.UniversalClient, action func(context.Context, *redis.Client) error) error {
	if c1, ok := client.(*redis.Client); ok {
		err := action(ctx, c1)
		return sderr.Wrapf(err, "for each shard error")
	} else if c1, ok := client.(*redis.Ring); ok {
		err := c1.ForEachShard(ctx, action)
		return sderr.Wrapf(err, "for each shard error (ring)")
	} else if c1, ok := client.(*redis.ClusterClient); ok {
		err := c1.ForEachShard(ctx, action)
		return sderr.Wrapf(err, "for each shard error (cluster)")
	} else {
		panic(sderr.Newf("for each shards error"))
	}
}

// IsNotFound 判断是否为redis的NotFound语义的错误
func IsNotFound(err error) bool {
	return sderr.Is(err, redis.Nil)
}
