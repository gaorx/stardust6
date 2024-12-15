package sdwebapp

import (
	"github.com/gaorx/stardust6/sderr"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"io/fs"
)

// App 描述一个Web应用
type App struct {
	*echo.Echo
}

// Options 描述App的选项
type Options struct {
	Debug          bool                  // 是否调试模式
	DisplayBanner  bool                  // 是否显示Banner
	DisplayPort    bool                  // 是否显示端口
	Binder         echo.Binder           // 绑定器
	Renderer       echo.Renderer         // 渲染器
	Fsys           fs.FS                 // 文件系统
	Validator      echo.Validator        // 验证器
	IPExtractor    echo.IPExtractor      // IP提取器
	ErrorHandler   echo.HTTPErrorHandler // 错误处理器
	SlogOptions    *SlogOptions          // slog选项
	EnableValidate bool                  // 是否启用验证
}

// New 创建一个App
func New(opts *Options) *App {
	opts1 := lo.FromPtr(opts)
	app := NewFrom(echo.New())
	app.Debug = opts1.Debug
	app.HideBanner = !opts1.DisplayBanner
	app.HidePort = !opts1.DisplayPort
	if opts1.Validator != nil {
		app.Validator = opts1.Validator
	} else {
		app.Validator = ValidatorFunc(defaultValidate)
	}
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
	app.Pre(
		SlogRecover(opts1.SlogOptions),
		func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				c.Set(akValidationEnabled, opts1.EnableValidate)
				return next(c)
			}
		},
	)
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
