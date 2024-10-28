package sdslogh

import (
	"io"
	"log/slog"
	"time"
)

type JsonOptions struct {
	AddSource   bool
	Level       slog.Leveler
	ReplaceAttr func(groups []string, a slog.Attr) slog.Attr

	Writer     io.Writer
	File       string
	TimeFormat string
	BufferSize int
}

func (b JsonOptions) Handler() slog.Handler {
	if b.TimeFormat == "" {
		b.TimeFormat = time.DateTime
	}
	w, _, err := newWriter(b.Writer, b.File, b.BufferSize)
	if err != nil {
		panic(err.Error())
	}
	return slog.NewJSONHandler(w, &slog.HandlerOptions{
		AddSource:   b.AddSource,
		Level:       b.Level,
		ReplaceAttr: setTimeFormat(b.TimeFormat, b.ReplaceAttr),
	})
}
