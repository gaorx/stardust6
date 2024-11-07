package sdslog

import (
	"context"
	"log/slog"
)

var (
	// DiscardHandler 不输出任何日志的Handler
	DiscardHandler Handler = discardHandler{}

	// DiscardMiddleware 不做任何处理的Middleware
	DiscardMiddleware Middleware = discardMiddleware

	// DiscardLogger 不输出任何日志的Logger
	DiscardLogger *Logger = slog.New(DiscardHandler)
)

type discardHandler struct{}

func (discardHandler) Enabled(context.Context, Level) bool {
	return false
}

func (discardHandler) Handle(context.Context, Record) error {
	return nil
}

func (h discardHandler) WithAttrs([]Attr) Handler {
	return h
}

func (h discardHandler) WithGroup(string) Handler {
	return h
}

func discardMiddleware(h Handler) Handler {
	return h
}
