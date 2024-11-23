package sdwebapp

import (
	"github.com/gaorx/stardust6/sderr"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
)

type App struct {
	*echo.Echo
}

type Options struct {
	DebugMode     bool
	DisplayBanner bool
}

func New() *App {
	return NewFrom(echo.New())
}

func NewFrom(e *echo.Echo) *App {
	return &App{Echo: e}
}

func (app *App) Install(components ...Component) error {
	for _, c := range components {
		if c != nil {
			err := c.Apply(app)
			if err != nil {
				return sderr.Newf("")
			}
		}
	}
	return nil
}

func (app *App) MustInstall(components ...Component) {
	lo.Must0(app.Install(components...))
}
