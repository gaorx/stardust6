package sdslogh

import (
	"log/slog"
)

type Builder interface {
	Handler() (slog.Handler, error)
}
