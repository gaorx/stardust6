package sdobjectstore

import (
	"github.com/gaorx/stardust6/sderr"
	"github.com/samber/lo"
	"strings"
)

// Config 对象存储配置
type Config struct {
	Type           string `json:"type" toml:"type"`
	Root           string `json:"root" toml:"root"`
	Endpoint       string `json:"endpoint" toml:"endpoint"`
	AccessKey      string `json:"access_key" toml:"access_key"`
	AccessSecret   string `json:"access_secret" toml:"access_secret"`
	Bucket         string `json:"bucket" toml:"bucket"`
	Prefix         string `json:"prefix" toml:"prefix"`
	InternalPrefix string `json:"internal_prefix" toml:"internal_prefix"`
}

// New 通过配置创建对象存储
func New(config *Config) (Store, error) {
	config1 := lo.FromPtr(config)
	switch strings.ToLower(config1.Type) {
	case "discard":
		return Store{Discard}, nil
	case "dir", "directory":
		return Store{dir{root: config1.Root}}, nil
	case "aliyun_oss", "aliyun-oss", "aliyunoss":
		aoss, err := NewAliyunOSS(&AliyunOssConfig{
			Endpoint:       config1.Endpoint,
			AccessKey:      config1.AccessKey,
			AccessSecret:   config1.AccessSecret,
			Bucket:         config1.Bucket,
			Prefix:         config1.Prefix,
			InternalPrefix: config1.InternalPrefix,
		})
		if err != nil {
			return Store{nil}, sderr.Wrap(err)
		}
		return Store{aoss}, nil
	default:
		return Store{nil}, sderr.With("type", config1.Type).Newf("illegal type")
	}
}
