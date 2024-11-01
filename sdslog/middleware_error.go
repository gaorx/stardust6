package sdslog

import (
	"fmt"
	"github.com/gaorx/stardust6/sderr"
	"github.com/samber/lo"
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

// ExpandErrorOptions 展开Error的Options
type ExpandErrorOptions struct {
	Stack bool // 是否打印错误的stacktrace
	Attrs bool // 是否打印sderr中的Attrs
}

func ExpandError(opts *ExpandErrorOptions) Middleware {
	opts1 := lo.FromPtr(opts)
	expand := func(a Value) Value {
		err, ok := a.Any().(error)
		if !ok || err == nil {
			return a
		}
		values := []Attr{slog.String("msg", err.Error())}
		if opts1.Attrs {
			attrs := sderr.Attrs(err)
			for k, v := range attrs {
				values = append(values, slog.String(k, fmt.Sprintf("%v", v)))
			}
		}
		if opts1.Stack {
			values = append(values, slog.String("stack", sderr.RootStack(err).Top().String()))
		}
		return slog.GroupValue(values...)
	}
	return func(h Handler) Handler {
		return Format(FormatByKey(ErrorKey, expand))(h)
	}
}
