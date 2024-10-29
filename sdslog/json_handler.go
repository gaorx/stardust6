package sdslog

import (
	"io"
	"log/slog"
	"time"
)

type JsonHandlerOptions struct {
	AddSource   bool
	Level       Leveler
	ReplaceAttr func(groups []string, a Attr) Attr

	Writer     io.Writer
	File       string
	TimeFormat string
	BufferSize int
}

func JsonFile(level Leveler, file string) Handler {
	return JsonHandlerOptions{
		Level: level,
		File:  file,
	}.NewHandler()
}

func NewJsonHandler(opts *JsonHandlerOptions) Handler {
	if opts == nil {
		panic("nil options")
	}
	return opts.NewHandler()
}

func (opts JsonHandlerOptions) NewHandler() Handler {
	if opts.TimeFormat == "" {
		opts.TimeFormat = time.DateTime
	}
	w, _, err := newWriter(opts.Writer, opts.File, opts.BufferSize)
	if err != nil {
		panic(err.Error())
	}
	return slog.NewJSONHandler(w, &slog.HandlerOptions{
		AddSource:   opts.AddSource,
		Level:       opts.Level,
		ReplaceAttr: setTimeFormat(opts.TimeFormat, opts.ReplaceAttr),
	})
}
