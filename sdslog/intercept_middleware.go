package sdslog

import (
	"context"
	slogmulti "github.com/samber/slog-multi"
)

type (
	InterceptEnabledFunc   = func(ctx context.Context, level Level, next func(context.Context, Level) bool) bool
	InterceptHandleFunc    = func(ctx context.Context, record Record, next func(context.Context, Record) error) error
	InterceptWithAttrsFunc = func(attrs []Attr, next func([]Attr) Handler) Handler
	InterceptWithGroupFunc = func(name string, next func(string) Handler) Handler
)

func InterceptAll(
	enabledFunc InterceptEnabledFunc,
	handleFunc InterceptHandleFunc,
	withAttrsFunc InterceptWithAttrsFunc,
	withGroupFunc InterceptWithGroupFunc,
) Middleware {
	if enabledFunc == nil && handleFunc == nil && withAttrsFunc == nil && withGroupFunc == nil {
		return DiscardMiddleware
	}
	if enabledFunc == nil {
		enabledFunc = func(ctx context.Context, level Level, next func(context.Context, Level) bool) bool {
			return next(ctx, level)
		}
	}
	if handleFunc == nil {
		handleFunc = func(ctx context.Context, record Record, next func(context.Context, Record) error) error {
			return next(ctx, record)
		}
	}
	if withAttrsFunc == nil {
		withAttrsFunc = func(attrs []Attr, next func([]Attr) Handler) Handler {
			return next(attrs)
		}
	}
	if withGroupFunc == nil {
		withGroupFunc = func(name string, next func(string) Handler) Handler {
			return next(name)
		}
	}
	return slogmulti.NewInlineMiddleware(enabledFunc, handleFunc, withAttrsFunc, withGroupFunc)
}

func InterceptEnabled(enabledFunc InterceptEnabledFunc) Middleware {
	if enabledFunc == nil {
		return DiscardMiddleware
	}
	return slogmulti.NewEnabledInlineMiddleware(enabledFunc)
}

func InterceptHandle(handleFunc InterceptHandleFunc) Middleware {
	if handleFunc == nil {
		return DiscardMiddleware
	}
	return slogmulti.NewHandleInlineMiddleware(handleFunc)
}

func InterceptWithAttrs(withAttrsFunc InterceptWithAttrsFunc) Middleware {
	if withAttrsFunc == nil {
		return DiscardMiddleware
	}
	return slogmulti.NewWithAttrsInlineMiddleware(withAttrsFunc)
}

func InterceptWithGroup(withGroupFunc InterceptWithGroupFunc) Middleware {
	if withGroupFunc == nil {
		return DiscardMiddleware
	}
	return slogmulti.NewWithGroupInlineMiddleware(withGroupFunc)
}
