package sdslog

import (
	"context"
	"log/slog"
)

var (
	DiscardHandler    Handler    = discardHandler{}         // 不输出任何日志的Handler
	DiscardMiddleware Middleware = discardMiddleware        // 不做任何处理的Middleware
	DiscardLogger     *Logger    = slog.New(DiscardHandler) // 不输出任何日志的Logger
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
