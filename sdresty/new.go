package sdresty

import (
	"crypto/tls"
	"github.com/go-resty/resty/v2"
	"github.com/samber/lo"
	"time"
)

// Options 创建resty.Client的选项
type Options struct {
	Timeout            time.Duration
	RetryCount         int
	Proxy              string
	InsecureSkipVerify bool
	QueryParams        map[string]any
	PathParams         map[string]any
	Headers            map[string]string
}

// New 创建resty.Client
func New(opts *Options) *resty.Client {
	opts1 := lo.FromPtr(opts)
	c := resty.New()
	c.SetTimeout(opts1.Timeout)
	c.SetRetryCount(opts1.RetryCount)
	if opts1.Proxy != "" {
		c.SetProxy(opts1.Proxy)
	} else {
		c.RemoveProxy()
	}
	if opts1.InsecureSkipVerify {
		c.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}
	if len(opts1.QueryParams) > 0 {
		c.SetQueryParams(toStringMap(opts1.QueryParams))
	}
	if len(opts1.PathParams) > 0 {
		c.SetPathParams(toStringMap(opts1.PathParams))
	}
	if len(opts1.Headers) > 0 {
		c.SetHeaders(opts1.Headers)
	}
	return c
}
