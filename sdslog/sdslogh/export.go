package sdslogh

import (
	slogmulti "github.com/samber/slog-multi"
	"log/slog"
)

func New(handlersOrBuildersOrMiddlewares ...any) slog.Handler {
	var handlers []slog.Handler
	var middlewares = []Middleware{expandError}

	addHandler := func(h slog.Handler) {
		if h != nil {
			handlers = append(handlers, h)
		}
	}

	for _, v := range handlersOrBuildersOrMiddlewares {
		if h, ok := v.(slog.Handler); ok && h != nil {
			handlers = append(handlers, h)
		} else if b, ok := v.(Builder); ok && b != nil {
			h, err := b.Handler()
			if err != nil {
				panic(err.Error())
			}
			addHandler(h)
		} else if b, ok := v.(interface{ NewHandler() (slog.Handler, error) }); ok && b != nil {
			h, err := b.NewHandler()
			if err != nil {
				panic(err.Error())
			}
			addHandler(h)
		} else if b, ok := v.(func() (slog.Handler, error)); ok && b != nil {
			h, err := b()
			if err != nil {
				panic(err.Error())
			}
			addHandler(h)
		} else if b, ok := v.(interface{ Handler() slog.Handler }); ok && b != nil {
			h := b.Handler()
			addHandler(h)
		} else if b, ok := v.(interface{ NewHandler() slog.Handler }); ok && b != nil {
			h := b.NewHandler()
			addHandler(h)
		} else if b, ok := v.(func() slog.Handler); ok && b != nil {
			h := b()
			addHandler(h)
		} else if file, ok := v.(string); ok {
			h, err := TextOptions{
				Level:  slog.LevelDebug,
				File:   file,
				Pretty: isStdoutFile(file) || isStderrFile(file),
			}.Handler()
			if err != nil {
				panic(err.Error())
			}
			addHandler(h)
		} else if m, ok := v.(Middleware); ok && m != nil {
			middlewares = append(middlewares, m)
		} else if m, ok := v.(func(slog.Handler) slog.Handler); ok && m != nil {
			middlewares = append(middlewares, m)
		}
	}
	switch len(handlers) {
	case 0:
		return Wrap(Discard, middlewares...)
	case 1:
		return Wrap(handlers[0], middlewares...)
	default:
		return Wrap(slogmulti.Fanout(handlers...), middlewares...)
	}
}
