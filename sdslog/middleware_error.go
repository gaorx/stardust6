package sdslog

import (
	"github.com/gaorx/stardust6/sderr"
	slogformatter "github.com/samber/slog-formatter"
	"log/slog"
)

const ErrorKey = "error" // 错误在slog的Attr中的key

// ErrorValue 将一个error转换为slog中的Attr，用于在日志中输出错误
func ErrorValue(err error) Value {
	return slog.AnyValue(err)
}

// E 将一个error转换为slog中的Attr(key的值为"error")，用于在日志中输出错误
func E(err error) Attr {
	return slog.Any(ErrorKey, err)
}

// ExpandSderrError 一个middleware，用户展开在slog中输出更详细的错误信息，对sderr.Error提供了特别支持，可以展开其中的Attrs作为slog中的attrs
func ExpandSderrError(h Handler) Handler {
	expandError := func(a Value) Value {
		err, ok := a.Any().(error)
		if !ok || err == nil {
			return a
		}
		values := []Attr{slog.String("msg", err.Error())}
		for k, v := range sderr.Attrs(err) {
			values = append(values, slog.Any(k, v))
		}
		return slog.GroupValue(values...)
	}
	return Format(FormatByKey(ErrorKey, expandError))(h)
}

// ExpandGenericError 一个middleware，用户展开在slog中输出更详细的的错误信息
func ExpandGenericError(h Handler) Handler {
	return Format(slogformatter.ErrorFormatter("error"))(h)
}
