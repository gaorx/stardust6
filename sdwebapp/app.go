package sdwebapp

import (
	"github.com/gaorx/stardust6/sderr"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"io/fs"
)

type App struct {
	*echo.Echo
}

type Options struct {
	Debug         bool
	DisplayBanner bool
	DisplayPort   bool
	Binder        echo.Binder
	Renderer      echo.Renderer
	Fsys          fs.FS
	IPExtractor   echo.IPExtractor
	ErrorHandler  echo.HTTPErrorHandler
	SlogOptions   *SlogOptions
}

func New(opts *Options) *App {
	opts1 := lo.FromPtr(opts)
	app := NewFrom(echo.New())
	app.Debug = opts1.Debug
	app.HideBanner = !opts1.DisplayBanner
	app.HidePort = !opts1.DisplayPort
	if opts1.Binder != nil {
		app.Binder = opts1.Binder
	}
	if opts1.Renderer != nil {
		app.Renderer = opts1.Renderer
	}
	if opts1.Fsys != nil {
		app.Filesystem = opts1.Fsys
	}
	if opts1.IPExtractor != nil {
		app.IPExtractor = opts1.IPExtractor
	}
	if opts1.ErrorHandler != nil {
		app.HTTPErrorHandler = opts1.ErrorHandler
	} else {
		app.HTTPErrorHandler = DefaultHttpErrorHandler
	}
	app.Use(SlogRecover(opts1.SlogOptions))
	return app
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
