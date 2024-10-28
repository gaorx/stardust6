package sdslog

import (
	"log/slog"
)

type Extension struct {
	Logger *slog.Logger
}

func X(l *slog.Logger) Extension {
	return Extension{Logger: l}
}
