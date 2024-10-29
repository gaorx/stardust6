package sdslog

import (
	"github.com/samber/lo"
	slogmulti "github.com/samber/slog-multi"
)

func Wrap(h Handler, middlewares ...Middleware) Handler {
	middlewares = lo.Filter(middlewares, func(m Middleware, _ int) bool {
		return m != nil
	})
	if len(middlewares) <= 0 {
		return h
	}
	return slogmulti.Pipe(middlewares...).Handler(h)
}
