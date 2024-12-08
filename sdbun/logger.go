package sdbun

import (
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/extra/bundebug"
	"log/slog"
)

// LoggerOf 创建一个bun.QueryHook
func LoggerOf(name string) bun.QueryHook {
	switch name {
	case "", "discard", "disable":
		return discardLogger{}
	case "default", "bun":
		return bundebug.NewQueryHook(bundebug.WithVerbose(true), bundebug.FromEnv("BUNDEBUG"))
	case "slog":
		return Slog(slog.Default())
	default:
		return discardLogger{}
	}
}
