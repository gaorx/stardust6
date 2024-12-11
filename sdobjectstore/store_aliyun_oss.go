package sdobjectstore

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/gaorx/stardust6/sderr"
	"github.com/samber/lo"
	"io"
	"strings"
)

// AliyunOssConfig 阿里云OSS配置
type AliyunOssConfig struct {
	Endpoint       string `json:"endpoint" toml:"endpoint"`
	AccessKey      string `json:"access_key" toml:"access_key"`
	AccessSecret   string `json:"access_secret" toml:"access_secret"`
	Bucket         string `json:"bucket" toml:"bucket"`
	Prefix         string `json:"prefix" toml:"prefix"`
	InternalPrefix string `json:"internal_prefix" toml:"internal_prefix"`
}

type aliyunOss struct {
	client         *oss.Client
	bucket         string
	prefix         string
	internalPrefix string
}

// NewAliyunOSS 通过配置创建阿里云OSS对象存储
func NewAliyunOSS(config *AliyunOssConfig) (Interface, error) {
	config1 := lo.FromPtr(config)
	if config1.Bucket == "" {
		return nil, sderr.Newf("no bucket")
	}
	if config1.Prefix == "" {
		return nil, sderr.Newf("no prefix")
	}
	if config1.InternalPrefix == "" {
		return nil, sderr.Newf("no internal prefix")
	}
	client, err := oss.New(config1.Endpoint, config1.AccessKey, config1.AccessSecret)
	if err != nil {
		return nil, sderr.Wrap(err)
	}
	return &aliyunOss{
		client:         client,
		bucket:         config1.Bucket,
		prefix:         config1.Prefix,
		internalPrefix: config1.InternalPrefix,
	}, nil
}

// Store 实现 Interface.Store
func (aoss *aliyunOss) Store(src Source, objectName string) (*Target, error) {
	if src == nil {
		return nil, sderr.Newf("nil source")
	}

	// 展开文件名称
	expandedObjectName, err := expandObjectName(src, objectName)
	if err != nil {
		return nil, err
	}

	// 针对特殊类型content-type，修正掉content-disposition
	ct, cd := src.ContentType(), "attachment"
	if strings.HasPrefix(ct, "text/") ||
		strings.HasPrefix(ct, "image/") ||
		strings.HasPrefix(ct, "video/") ||
		strings.HasPrefix(ct, "audio/") ||
		strings.Contains(ct, "json") ||
		strings.Contains(ct, "javascript") ||
		strings.Contains(ct, "ecmascript") {
		cd = "inline" // 阿里云实际上会忽略这个inline，强制改为attachment
	}

	// 存储
	b, err := aoss.client.Bucket(aoss.bucket)
	if err != nil {
		return nil, sderr.Wrap(err)
	}
	r, err := src.Open()
	if err != nil {
		return nil, sderr.Wrap(err)
	}
	defer func() {
		if rc, ok := r.(io.ReadCloser); ok {
			_ = rc.Close()
		}
	}()
	err = b.PutObject(expandedObjectName, r, oss.ContentType(src.ContentType()), oss.ContentDisposition(cd))
	if err != nil {
		return nil, sderr.Wrap(err)
	}
	// 返回信息
	return &Target{
		Typ:            HttpTarget,
		Prefix:         aoss.prefix,
		InternalPrefix: aoss.internalPrefix,
		Path:           expandedObjectName,
	}, nil
}
