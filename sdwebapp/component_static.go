package sdwebapp

import (
	"github.com/labstack/echo/v4"
	"io/fs"
	"net/http"
	"path"
)

type StaticContent struct {
	Method            string
	Path              string
	StatusCode        int
	ContentType       string
	Content           []byte
	Guard             RouteGuard
	GuardErrorHandler echo.HTTPErrorHandler
	Middlewares       []MiddlewareFunc
}

type StaticFile struct {
	Path              string
	Fsys              fs.FS
	Filename          string
	Guard             RouteGuard
	GuardErrorHandler echo.HTTPErrorHandler
	Middlewares       []MiddlewareFunc
}

type StaticDirectory struct {
	Prefix                string
	Dirname               string
	Fsys                  fs.FS
	DisablePathUnescaping bool
	RedirectToSlash       bool
	Fallback              func(p string) []string
	Guard                 RouteGuard
	GuardErrorHandler     echo.HTTPErrorHandler
	Middlewares           []MiddlewareFunc
}

var (
	_ Routable  = StaticContent{}
	_ Component = StaticContent{}
	_ Routable  = StaticFile{}
	_ Component = StaticFile{}
	_ Routable  = StaticDirectory{}
	_ Component = StaticDirectory{}
)

func Blob(path string, contentType string, content []byte, middlewares ...MiddlewareFunc) *StaticContent {
	return &StaticContent{Path: path, ContentType: contentType, Content: content, Middlewares: middlewares}
}

func Text(path string, contentType string, content string, middlewares ...MiddlewareFunc) *StaticContent {
	if contentType == "" {
		contentType = "text/plain"
	}
	return Blob(path, contentType, []byte(content), middlewares...)
}

func HTML(path string, content string, middlewares ...MiddlewareFunc) *StaticContent {
	return Text(path, "text/html", content, middlewares...)
}

func Javascript(path string, content string, middlewares ...MiddlewareFunc) *StaticContent {
	return Text(path, "application/javascript", content, middlewares...)
}

func CSS(path string, content string, middlewares ...MiddlewareFunc) *StaticContent {
	return Text(path, "text/css", content, middlewares...)
}

func (c *StaticContent) SetMethod(method string) *StaticContent {
	c.Method = method
	return c
}

func (c *StaticContent) SetStatusCode(statusCode int) *StaticContent {
	c.StatusCode = statusCode
	return c
}

func (c *StaticContent) SetContentType(contentType string) *StaticContent {
	c.ContentType = contentType
	return c
}

func (c *StaticContent) SetGuard(guard RouteGuard) *StaticContent {
	c.Guard = guard
	return c
}

func (c *StaticContent) SetGuardErrorHandler(handler echo.HTTPErrorHandler) *StaticContent {
	c.GuardErrorHandler = handler
	return c
}

func (c *StaticContent) AddMiddleware(middlewares ...MiddlewareFunc) *StaticContent {
	c.Middlewares = append(c.Middlewares, middlewares...)
	return c
}

func (c StaticContent) ToRoutes(*App) []*Route {
	c1 := c
	if c1.Method == "" {
		c1.Method = http.MethodGet
	}
	if c1.StatusCode == 0 {
		c1.StatusCode = http.StatusOK
	}
	if c1.ContentType == "" {
		c1.ContentType = "application/octet-stream"
	}
	guard := c.Guard
	if guard == nil {
		guard = PermitAll()
	}
	guardErrorHandler := c.GuardErrorHandler
	if guardErrorHandler == nil {
		guardErrorHandler = defaultRouteErrorHandler
	}
	return []*Route{{
		Method: http.MethodGet,
		Path:   c.Path,
		Handler: func(c echo.Context) error {
			return c.Blob(c1.StatusCode, c1.ContentType, c1.Content)
		},
		Guard:             guard,
		GuardErrorHandler: guardErrorHandler,
		Middlewares:       c.Middlewares,
	}}
}

func (c StaticContent) Apply(app *App) error {
	return applyRoutes(app, c.ToRoutes(app))
}

func File(path string, filename string, middlewares ...MiddlewareFunc) *StaticFile {
	return &StaticFile{Path: path, Filename: filename, Middlewares: middlewares}
}

func FileFS(path string, filename string, fsys fs.FS, middlewares ...MiddlewareFunc) *StaticFile {
	return &StaticFile{Path: path, Filename: filename, Fsys: fsys, Middlewares: middlewares}
}

func (f *StaticFile) SetGuard(guard RouteGuard) *StaticFile {
	f.Guard = guard
	return f
}

func (f *StaticFile) SetGuardErrorHandler(handler echo.HTTPErrorHandler) *StaticFile {
	f.GuardErrorHandler = handler
	return f
}

func (f *StaticFile) AddMiddleware(middlewares ...MiddlewareFunc) *StaticFile {
	f.Middlewares = append(f.Middlewares, middlewares...)
	return f
}

func (f StaticFile) ToRoutes(*App) []*Route {
	guard := f.Guard
	if guard == nil {
		guard = PermitAll()
	}
	guardErrorHandler := f.GuardErrorHandler
	if guardErrorHandler == nil {
		guardErrorHandler = defaultRouteErrorHandler
	}
	route := Route{
		Method:            http.MethodGet,
		Path:              f.Path,
		Guard:             guard,
		GuardErrorHandler: guardErrorHandler,
		Middlewares:       f.Middlewares,
	}
	if f.Fsys != nil {
		route.Handler = echo.StaticFileHandler(f.Filename, f.Fsys)
	} else {
		route.Handler = func(c echo.Context) error {
			return c.File(f.Filename)
		}
	}
	return []*Route{&route}
}

func (f StaticFile) Apply(app *App) error {
	return applyRoutes(app, f.ToRoutes(app))
}

func Dir(prefix, dirname string, middlewares ...MiddlewareFunc) *StaticDirectory {
	return &StaticDirectory{Prefix: prefix, Dirname: dirname, Middlewares: middlewares}
}

func DirFS(prefix string, fsys fs.FS, middlewares ...MiddlewareFunc) *StaticDirectory {
	return &StaticDirectory{Prefix: prefix, Fsys: fsys, Middlewares: middlewares}
}

func (d *StaticDirectory) SetDisablePathUnescaping(b bool) *StaticDirectory {
	d.DisablePathUnescaping = b
	return d
}

func (d *StaticDirectory) SetRedirectToSlash(b bool) *StaticDirectory {
	d.RedirectToSlash = b
	return d
}

func (d *StaticDirectory) SetFallback(f func(p string) []string) *StaticDirectory {
	d.Fallback = f
	return d
}

func (d *StaticDirectory) SetGuard(guard RouteGuard) *StaticDirectory {
	d.Guard = guard
	return d
}

func (d *StaticDirectory) SetGuardErrorHandler(handler echo.HTTPErrorHandler) *StaticDirectory {
	d.GuardErrorHandler = handler
	return d
}

func (d *StaticDirectory) AddMiddlewares(middlewares ...MiddlewareFunc) *StaticDirectory {
	d.Middlewares = append(d.Middlewares, middlewares...)
	return d
}

func (d *StaticDirectory) SetFallbackToRootIndexPage() *StaticDirectory {
	d.Fallback = func(p string) []string {
		ext := path.Ext(p)
		if ext == "" || ext == ".html" || ext == ".htm" {
			// 仅仅针对想访问页面的路径才回尝试返回跟目录下的 index.html，对css和js等资源文件都不会尝试
			return []string{"/" + indexPage}
		}
		return nil
	}
	return d
}

func (d StaticDirectory) ToRoutes(app *App) []*Route {
	var fsys fs.FS
	if d.Fsys != nil {
		fsys = d.Fsys
	} else {
		fsys = echo.MustSubFS(app.Echo.Filesystem, d.Dirname)
	}
	guard := d.Guard
	if guard == nil {
		guard = PermitAll()
	}
	guardErrorHandler := d.GuardErrorHandler
	if guardErrorHandler == nil {
		guardErrorHandler = defaultRouteErrorHandler
	}
	return []*Route{{
		Method: http.MethodGet,
		Path:   d.Prefix + "*",
		Handler: StaticDirectoryHandler(fsys, &StaticDirectoryOptions{
			DisablePathUnescaping: d.DisablePathUnescaping,
			RedirectToSlash:       d.RedirectToSlash,
			Fallback:              d.Fallback,
		}),
		Guard:             guard,
		GuardErrorHandler: guardErrorHandler,
		Middlewares:       d.Middlewares,
	}}
}

func (d StaticDirectory) Apply(app *App) error {
	return applyRoutes(app, d.ToRoutes(app))
}
