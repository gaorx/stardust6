package sdslog

import (
	"github.com/lmittmann/tint"
	"io"
	"log/slog"
	"time"
)

// TextHandlerOptions 生成Json handler的选项
type TextHandlerOptions struct {
	AddSource   bool                               // 日志中是否加入源码位置信息
	Level       Leveler                            // 日志级别
	ReplaceAttr func(groups []string, a Attr) Attr // 替换Record中Attr的函数

	Writer     io.Writer // 指定的输出流，如果为nil，则使用File指向的文件
	File       string    // 输出的文件名，为stdout时输出到标准输出，为stderr时输出到标准错误，为discard时丢弃日志
	TimeFormat string    // 日志输出的时间格式，默认为time.DateTime
	BufferSize int       // 日志缓冲区大小，默认为0，没有缓冲区，如果这个值为4096，则表示有一个4KB的缓冲区

	Pretty bool // 是否输出更便于阅读的格式
}

// TextFile 构造一个输出到文件的text handler，level为日志级别，file为输出文件名, pretty为true则输出更便于阅读的格式
func TextFile(level Leveler, file string, pretty bool) Handler {
	return TextHandlerOptions{
		Level:  level,
		File:   file,
		Pretty: pretty,
	}.NewHandler()
}

// NewTextHandler 根据选项生成一个text handler
func NewTextHandler(opts *TextHandlerOptions) Handler {
	if opts == nil {
		panic("nil options")
	}
	return opts.NewHandler()
}

// NewHandler 根据此选项生成一个text handler
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
