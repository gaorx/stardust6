package sdslogh

import (
	"github.com/samber/lo"
	slogmulti "github.com/samber/slog-multi"
	"log/slog"
)

type Middleware = slogmulti.Middleware

func Wrap(h slog.Handler, middlewares ...Middleware) slog.Handler {
	middlewares = lo.Filter(middlewares, func(m Middleware, _ int) bool {
		return m != nil
	})
	if len(middlewares) <= 0 {
		return h
	}
	return slogmulti.Pipe(middlewares...).Handler(h)
}
