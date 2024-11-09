package sdresty

import (
	"github.com/go-resty/resty/v2"
)

// Extractor 从resty.Response中提取最终要获取的数据
type Extractor[T any] interface {
	// Extract 从resty.Response中提取数据，err参数就是Execute返回的错误值
	Extract(resp *resty.Response, err error) (T, error)

	// Func 从一个callback结果中提取resty.Response中数据，err参数就是Execute返回的错误值
	Func(f func() (*resty.Response, error)) (T, error)
}

// ExtractorFunc 从resty.Response中提取数据的函数
type ExtractorFunc[T any] func(*resty.Response, error) (T, error)

func (f ExtractorFunc[T]) Extract(resp *resty.Response, err error) (T, error) {
	return f(resp, err)
}

func (f ExtractorFunc[T]) Func(ff func() (*resty.Response, error)) (T, error) {
	resp, err := ff()
	return f(resp, err)
}
