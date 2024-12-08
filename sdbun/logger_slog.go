package sdbun

import (
	"context"
	"github.com/gaorx/stardust6/sdslog"
	"github.com/uptrace/bun"
	"log/slog"
	"time"
)

// Slog 创建一个使用slog实现的bun.QueryHook
func Slog(logger *slog.Logger) bun.QueryHook {
	return slogLogger{logger: logger}
}

type slogLogger struct {
	logger *slog.Logger
}

func (h slogLogger) BeforeQuery(ctx context.Context, e *bun.QueryEvent) context.Context {
	return ctx
}

func (h slogLogger) AfterQuery(ctx context.Context, e *bun.QueryEvent) {
	elapsed := time.Now().Sub(e.StartTime).String()
	if e.Err != nil {
		slog.With(
			sdslog.E(e.Err),
			slog.String("q", e.Query),
			slog.String("elapsed", elapsed),
			slog.String("op", e.Operation()),
		).Error("query error")
	} else {
		slog.With(
			slog.String("q", e.Query),
			slog.String("elapsed", elapsed),
			slog.String("op", e.Operation()),
		).Debug("query done")
	}
}
