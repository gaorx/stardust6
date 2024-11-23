package sdwebapp

import (
	"github.com/labstack/echo/v4"
)

type (
	HandlerFunc    = echo.HandlerFunc
	MiddlewareFunc = echo.MiddlewareFunc
)

type Component interface {
	Apply(*App) error
}

type ComponentFunc func(*App) error

func (f ComponentFunc) Apply(app *App) error {
	if f == nil {
		return nil
	}
	return f(app)
}
