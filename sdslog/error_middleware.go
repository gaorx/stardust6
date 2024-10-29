package sdslog

import (
	"github.com/gaorx/stardust6/sderr"
	slogformatter "github.com/samber/slog-formatter"
	"log/slog"
)

const ErrorKey = "error"

func ErrorValue(err error) Value {
	return slog.AnyValue(err)
}

func E(err error) Attr {
	return slog.Any(ErrorKey, err)
}

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

func ExpandGenericError(h Handler) Handler {
	return Format(slogformatter.ErrorFormatter("error"))(h)
}
