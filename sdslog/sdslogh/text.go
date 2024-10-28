package sdslogh

import (
	"github.com/gaorx/stardust6/sderr"
	"github.com/lmittmann/tint"
	"io"
	"log/slog"
	"time"
)

type TextOptions struct {
	AddSource   bool
	Level       slog.Leveler
	ReplaceAttr func(groups []string, a slog.Attr) slog.Attr

	Writer     io.Writer
	File       string
	TimeFormat string
	Pretty     bool
	BufferSize int
}

var _ Builder = TextOptions{}

func (b TextOptions) Handler() (slog.Handler, error) {
	if b.TimeFormat == "" {
		b.TimeFormat = time.DateTime
	}
	w, supportColored, err := newWriter(b.Writer, b.File, b.BufferSize)
	if err != nil {
		return nil, sderr.Wrap(err)
	}
	if b.Pretty {
		return tint.NewHandler(w, &tint.Options{
			AddSource:   b.AddSource,
			Level:       b.Level,
			ReplaceAttr: b.ReplaceAttr,
			TimeFormat:  b.TimeFormat,
			NoColor:     !supportColored,
		}), nil
	} else {
		return slog.NewTextHandler(w, &slog.HandlerOptions{
			AddSource:   b.AddSource,
			Level:       b.Level,
			ReplaceAttr: setTimeFormat(b.TimeFormat, b.ReplaceAttr),
		}), nil
	}
}
