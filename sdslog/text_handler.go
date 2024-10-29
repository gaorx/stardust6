package sdslog

import (
	"github.com/lmittmann/tint"
	"io"
	"log/slog"
	"time"
)

type TextHandlerOptions struct {
	AddSource   bool
	Level       Leveler
	ReplaceAttr func(groups []string, a Attr) Attr

	Writer     io.Writer
	File       string
	TimeFormat string
	Pretty     bool
	BufferSize int
}

func TextFile(level Leveler, file string) Handler {
	return TextHandlerOptions{
		Level:  level,
		File:   file,
		Pretty: (isStdoutFile(file) && isStdoutTerm()) || (isStderrFile(file) && isStderrTerm()),
	}.NewHandler()
}

func NewTextHandler(opts *TextHandlerOptions) Handler {
	if opts == nil {
		panic("nil options")
	}
	return opts.NewHandler()
}

func (opts TextHandlerOptions) NewHandler() Handler {
	if opts.TimeFormat == "" {
		opts.TimeFormat = time.DateTime
	}
	w, supportColored, err := newWriter(opts.Writer, opts.File, opts.BufferSize)
	if err != nil {
		panic(err.Error())
	}
	if opts.Pretty {
		return tint.NewHandler(w, &tint.Options{
			AddSource:   opts.AddSource,
			Level:       opts.Level,
			ReplaceAttr: opts.ReplaceAttr,
			TimeFormat:  opts.TimeFormat,
			NoColor:     !supportColored,
		})
	} else {
		return slog.NewTextHandler(w, &slog.HandlerOptions{
			AddSource:   opts.AddSource,
			Level:       opts.Level,
			ReplaceAttr: setTimeFormat(opts.TimeFormat, opts.ReplaceAttr),
		})
	}
}
