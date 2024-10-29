package sdslog

import (
	"context"
	"log/slog"
)

var (
	DiscardHandler    Handler    = discardHandler{}
	DiscardMiddleware Middleware = discardMiddleware
	DiscardLogger     *Logger    = slog.New(DiscardHandler)
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
