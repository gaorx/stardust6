package sdslog

import (
	"context"
	slogmock "github.com/samber/slog-mock"
)

// HandleRecord 通过一个Record处理函数来生成一个slog的Handler，用于自定义日志输出
func HandleRecord(h func(ctx context.Context, record Record) error) Handler {
	if h == nil {
		return DiscardHandler
	}
	return slogmock.Option{
		Handle: h,
	}.NewMockHandler()
}
