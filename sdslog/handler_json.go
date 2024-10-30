package sdslog

import (
	"io"
	"log/slog"
	"time"
)

// JsonHandlerOptions 生成Json handler的选项
type JsonHandlerOptions struct {
	AddSource   bool                               // 日志中是否加入源码位置信息
	Level       Leveler                            // 日志级别
	ReplaceAttr func(groups []string, a Attr) Attr // 替换Record中Attr的函数

	Writer     io.Writer // 指定的输出流，如果为nil，则使用File指向的文件
	File       string    // 输出的文件名，为stdout时输出到标准输出，为stderr时输出到标准错误，为discard时丢弃日志
	TimeFormat string    // 日志输出的时间格式，默认为time.DateTime
	BufferSize int       // 日志缓冲区大小，默认为0，没有缓冲区，如果这个值为4096，则表示有一个4KB的缓冲区
}

// JsonFile 构造一个输出到文件的json handler，level为日志级别，file为输出文件名
func JsonFile(level Leveler, file string) Handler {
	return JsonHandlerOptions{
		Level: level,
		File:  file,
	}.NewHandler()
}

// NewJsonHandler 根据选项生成一个json handler
func NewJsonHandler(opts *JsonHandlerOptions) Handler {
	if opts == nil {
		panic("nil options")
	}
	return opts.NewHandler()
}

// NewHandler 根据此选项生成一个json handler
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
