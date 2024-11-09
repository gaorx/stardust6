package sdresty

import (
	"context"
	"github.com/go-resty/resty/v2"
	"net/http"
)

// Execute 从resty.Client执行请求，并且提取最终要获取的数据
func Execute[T any](ctx context.Context, c *resty.Client, method, url string, executor Extractor[T], opts ...RequestOption) (T, error) {
	req := c.R().SetContext(ctx)
	for _, opt := range opts {
		if opt != nil {
			req = opt(req)
		}
	}
	resp, err := req.Execute(method, url)
	return executor.Extract(resp, err)
}

// GET 从resty.Client执行GET请求，并且提取最终要获取的数据
func GET[T any](ctx context.Context, c *resty.Client, url string, executor Extractor[T], opts ...RequestOption) (T, error) {
	return Execute(ctx, c, http.MethodGet, url, executor, opts...)
}

// POST 从resty.Client执行POST请求，并且提取最终要获取的数据
func POST[T any](ctx context.Context, c *resty.Client, url string, executor Extractor[T], data RequestOption, opts ...RequestOption) (T, error) {
	merged := append([]RequestOption{data}, opts...)
	return Execute(ctx, c, http.MethodPost, url, executor, merged...)
}

// PUT 从resty.Client执行PUT请求，并且提取最终要获取的数据
func PUT[T any](ctx context.Context, c *resty.Client, url string, executor Extractor[T], data RequestOption, opts ...RequestOption) (T, error) {
	merged := append([]RequestOption{data}, opts...)
	return Execute(ctx, c, http.MethodPut, url, executor, merged...)
}

// PATCH 从resty.Client执行PATCH请求，并且提取最终要获取的数据
func PATCH[T any](ctx context.Context, c *resty.Client, url string, executor Extractor[T], data RequestOption, opts ...RequestOption) (T, error) {
	merged := append([]RequestOption{data}, opts...)
	return Execute(ctx, c, http.MethodPatch, url, executor, merged...)
}

// DELETE 从resty.Client执行DELETE请求，并且提取最终要获取的数据
func DELETE[T any](ctx context.Context, c *resty.Client, url string, executor Extractor[T], opts ...RequestOption) (T, error) {
	return Execute(ctx, c, http.MethodDelete, url, executor, opts...)
}

// CURL 从resty.Client执行GET，并提取body中的文本
func CURL(ctx context.Context, c *resty.Client, url string, opts ...RequestOption) (string, error) {
	return GET(ctx, c, url, ForTextBody(), opts...)
}
