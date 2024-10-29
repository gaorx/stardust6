package sdslog

import (
	"context"
)

var DiscardHandler Handler = discardHandler{}

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

var DiscardMiddleware Middleware = func(h Handler) Handler {
	return h
}
