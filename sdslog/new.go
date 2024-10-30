package sdslog

import (
	"github.com/samber/lo"
	slogmulti "github.com/samber/slog-multi"
	"log/slog"
)

// New 通过多个handler和多个middleware构建一个logger，如果有多个handler则会同时输出多路日志
func New(handlers []Handler, middlewares []Middleware) *slog.Logger {
	return slog.New(NewHandler(handlers, middlewares))
}

// NewHandler 通过多个handler和多个middleware构建一个handler，如果有多个handler则会同时输出多路日志
func NewHandler(handlers []Handler, middlewares []Middleware) Handler {
	handlers = lo.Filter(handlers, func(h Handler, _ int) bool {
		return h != nil
	})
	middlewares = lo.Filter(middlewares, func(m Middleware, _ int) bool {
		return m != nil
	})
	switch len(handlers) {
	case 0:
		return Wrap(DiscardHandler, middlewares...)
	case 1:
		return Wrap(handlers[0], middlewares...)
	default:
		return Wrap(slogmulti.Fanout(handlers...), middlewares...)
	}
}

// SetDefault 通过多个handler和多个middleware构建一个logger，并设置为全局默认的logger
func SetDefault(handlers []Handler, middlewares []Middleware) {
	slog.SetDefault(New(handlers, middlewares))
}
