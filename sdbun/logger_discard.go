package sdbun

import (
	"context"
	"github.com/uptrace/bun"
)

type discardLogger struct{}

func (h discardLogger) BeforeQuery(ctx context.Context, e *bun.QueryEvent) context.Context {
	return ctx
}

func (h discardLogger) AfterQuery(ctx context.Context, e *bun.QueryEvent) {
}
