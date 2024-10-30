package sdslog

import (
	"github.com/samber/lo"
	slogformatter "github.com/samber/slog-formatter"
	"time"
)

type Formatter = slogformatter.Formatter

// Format 构造一个用于格式化Record中的Attr值的middleware，支持多个Formatter来构造，
// 而formatter则是通过下面的FormatKey...来构建
func Format(formatters ...Formatter) Middleware {
	formatters = lo.Filter(formatters, func(f Formatter, _ int) bool {
		return f != nil
	})
	if len(formatters) <= 0 {
		return DiscardMiddleware
	}
	return slogformatter.NewFormatterMiddleware(formatters...)
}

func FormatByKind(kind Kind, formatter func(Value) Value) Formatter {
	return slogformatter.FormatByKind(kind, formatter)
}

func FormatByKey(key string, formatter func(Value) Value) Formatter {
	return slogformatter.FormatByKey(key, formatter)
}

func FormatByType[T any](formatter func(T) Value) Formatter {
	return slogformatter.FormatByType(formatter)
}

func FormatByFieldType[T any](key string, formatter func(T) Value) Formatter {
	return slogformatter.FormatByFieldType(key, formatter)
}

func FormatByGroup(groups []string, formatter func([]Attr) Value) Formatter {
	return slogformatter.FormatByGroup(groups, formatter)
}

func FormatByGroupKey(groups []string, key string, formatter func(Value) Value) Formatter {
	return slogformatter.FormatByGroupKey(groups, key, formatter)
}

func FormatByGroupKeyType[T any](groups []string, key string, formatter func(T) Value) Formatter {
	return slogformatter.FormatByGroupKeyType(groups, key, formatter)
}

func FormatTime(format string, location *time.Location) Formatter {
	return slogformatter.TimeFormatter(format, location)
}

func FormatUnixTimestamp(precision time.Duration) Formatter {
	return slogformatter.UnixTimestampFormatter(precision)
}

func ConvertTimezone(location *time.Location) Formatter {
	return slogformatter.TimezoneConverter(location)
}
