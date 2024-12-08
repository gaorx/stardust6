package sdredis

import (
	"github.com/gaorx/stardust6/sderr"
	"github.com/redis/go-redis/v9"
	"time"
)

// Address Redis地址
type Address struct {
	Addrs      []string `json:"addrs" toml:"addrs" yaml:"addrs"`
	DB         int      `json:"db" toml:"db" yaml:"db"`
	Password   string   `json:"password" toml:"password" yaml:"password"`
	MaxRetries int      `json:"max_retries" toml:"max_retries" yaml:"max_retries"`
	Cluster    bool     `json:"cluster" toml:"cluster" yaml:"cluster"`
}

// Dial 通过地址连接到Redis server
func Dial(addr Address) (redis.UniversalClient, error) {
	const (
		defaultPoolSize    = 30
		defaultPoolTimeout = 60 * time.Second
	)

	switch len(addr.Addrs) {
	case 0:
		return nil, sderr.Newf("no addresses")
	case 1:
		client := redis.NewClient(&redis.Options{
			Addr:        addr.Addrs[0],
			Password:    addr.Password,
			DB:          addr.DB,
			MaxRetries:  addr.MaxRetries,
			PoolSize:    defaultPoolSize,
			PoolTimeout: defaultPoolTimeout,
		})
		return client, nil
	default:
		if addr.Cluster {
			client := redis.NewClusterClient(&redis.ClusterOptions{
				Addrs:       addr.Addrs,
				Password:    addr.Password,
				MaxRetries:  addr.MaxRetries,
				PoolSize:    defaultPoolSize,
				PoolTimeout: defaultPoolTimeout,
			})
			return client, nil
		} else {
			addrMap := map[string]string{}
			for _, addr1 := range addr.Addrs {
				addrMap[addr1] = addr1
			}
			client := redis.NewRing(&redis.RingOptions{
				Addrs:       addrMap,
				Password:    addr.Password,
				DB:          addr.DB,
				MaxRetries:  addr.MaxRetries,
				PoolSize:    defaultPoolSize,
				PoolTimeout: defaultPoolTimeout,
			})
			return client, nil
		}
	}
}
