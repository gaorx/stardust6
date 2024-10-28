package sdslogh

import (
	"bufio"
	"fmt"
	"github.com/gaorx/stardust6/sderr"
	"github.com/mattn/go-isatty"
	slogformatter "github.com/samber/slog-formatter"
	"io"
	"log/slog"
	"os"
	"strings"
)

func isStdoutFile(file string) bool {
	file = strings.ToLower(file)
	return file == "" || file == "stdout"
}

func isStderrFile(file string) bool {
	file = strings.ToLower(file)
	return file == "stderr"
}

func isDiscardFile(file string) bool {
	file = strings.ToLower(file)
	return file == "discard"
}

func isTerm(fd uintptr) bool {
	return isatty.IsTerminal(fd) || isatty.IsCygwinTerminal(fd)
}

func isStdoutTerm() bool {
	fd := os.Stdout.Fd()
	return isTerm(fd)
}

func isStderrTerm() bool {
	fd := os.Stderr.Fd()
	return isTerm(fd)
}

func newWriter(w io.Writer, file string, bufferSize int) (io.Writer, bool, error) {
	if w != nil {
		if f, ok := w.(*os.File); ok && f != nil {
			fd := f.Fd()
			return w, isTerm(fd), nil
		} else {
			return w, false, nil
		}
	}
	if isStdoutFile(file) {
		return os.Stdout, isStdoutTerm(), nil
	} else if isStderrFile(file) {
		return os.Stderr, isStderrTerm(), nil
	} else if isDiscardFile(file) {
		return io.Discard, false, nil
	} else {
		var f io.Writer
		f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			return nil, false, sderr.With("file", file).Wrapf(err, "open log file error")
		}
		if bufferSize > 0 {
			f = bufio.NewWriterSize(w, bufferSize)
		}
		return f, false, nil
	}
}

func setTimeFormat(format string, replaceAttr func(groups []string, a slog.Attr) slog.Attr) func(groups []string, a slog.Attr) slog.Attr {
	return func(groups []string, a slog.Attr) slog.Attr {
		if replaceAttr != nil {
			a = replaceAttr(groups, a)
		}
		if a.Key != slog.TimeKey || a.Value.Kind() != slog.KindTime {
			return a
		}
		t := a.Value.Time()
		a.Value = slog.StringValue(t.Format(format))
		return a
	}
}

func expandError(h slog.Handler) slog.Handler {
	expandErr := func(err error) slog.Value {
		if err == nil {
			return slog.StringValue("<nil>")
		}
		values := []slog.Attr{slog.String("msg", err.Error())}
		for k, v := range sderr.Attrs(err) {
			values = append(values, slog.String(k, fmt.Sprintf("%v", v)))
		}
		// values = append(values, slog.("rootstack", "-a\n-b"))
		return slog.GroupValue(values...)
	}
	return slogformatter.NewFormatterHandler(
		slogformatter.FormatByFieldType("error", expandErr),
	)(h)
}
