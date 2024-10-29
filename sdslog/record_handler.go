package sdslog

import (
	"context"
	slogmock "github.com/samber/slog-mock"
)

func HandleRecord(h func(ctx context.Context, record Record) error) Handler {
	if h == nil {
		return DiscardHandler
	}
	return slogmock.Option{
		Handle: h,
	}.NewMockHandler()
}
