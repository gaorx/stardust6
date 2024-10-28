package sdslog

import (
	"github.com/gaorx/stardust6/sdslog/sdslogh"
	"log/slog"
)

func New(handlerOrBuilders ...any) *slog.Logger {
	return slog.New(sdslogh.New(handlerOrBuilders...))
}

func SetDefault(handlerOrBuilders ...any) {
	slog.SetDefault(New(handlerOrBuilders...))
}
