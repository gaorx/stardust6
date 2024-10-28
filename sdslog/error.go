package sdslog

import (
	"context"
	"fmt"
	"github.com/gaorx/stardust6/sderr"
	"log/slog"
)

func (ext Extension) ErrorContext(ctx context.Context, msg string, err error) {
	if err == nil {
		return
	}
	if msg != "" {
		msg = msg + ": " + err.Error()
	} else {
		msg = err.Error()
	}
	errAttrs := sderr.Attrs(err)
	var logAttrs []slog.Attr
	for k, v := range errAttrs {
		if k == "" {
			continue
		}
		logAttrs = append(logAttrs, slog.Any(k, v))
	}
	if false { // 这里可以打印root stack的栈顶frame
		rootStackFrames := sderr.RootStack(err).Frames()
		if len(rootStackFrames) > 0 {
			topFrame := rootStackFrames[0]
			logAttrs = append(logAttrs, slog.Attr{
				Key:   "stacktrace",
				Value: slog.StringValue(fmt.Sprintf("%s:%d", topFrame.File, topFrame.Line)),
			})
		}
	}

	ext.Logger.LogAttrs(ctx, slog.LevelError, msg, logAttrs...)
}

func (ext Extension) Error(msg string, err error) {
	ext.ErrorContext(context.Background(), msg, err)
}

func ErrorContext(ctx context.Context, msg string, err error) {
	X(slog.Default()).ErrorContext(ctx, msg, err)
}

func Error(msg string, err error) {
	X(slog.Default()).ErrorContext(context.Background(), msg, err)
}
