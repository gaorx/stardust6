package sdslog

import (
	slogmulti "github.com/samber/slog-multi"
	"log/slog"
)

type (
	Logger     = slog.Logger
	Handler    = slog.Handler
	Middleware = slogmulti.Middleware
	Record     = slog.Record
	Value      = slog.Value
	Attr       = slog.Attr
	Level      = slog.Level
	Leveler    = slog.Leveler
	Kind       = slog.Kind
)
