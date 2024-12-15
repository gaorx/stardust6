package sdsimpleapi

import (
	"github.com/gaorx/stardust6/sderr"
	"github.com/gaorx/stardust6/sdwebapp"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"net/http"
)

// API 描述一个API，它是一个Routable，也是一个Component
type API struct {
	Name              string
	Path              string
	Handler           any
	Middlewares       []echo.MiddlewareFunc
	Guard             sdwebapp.RouteGuard
	GuardErrorHandler echo.HTTPErrorHandler
	sdwebapp.Doc
}

var _ sdwebapp.Routable = (*API)(nil)

// R 创建一个API
func R(path string, handler any, guard sdwebapp.RouteGuard) *API {
	return &API{Path: path, Handler: handler, Guard: guard}
}

// SetName 设置API的名称
func (a *API) SetName(name string) *API {
	a.Name = name
	return a
}

// SetDoc 设置API的文档
func (a *API) SetDoc(doc *sdwebapp.Doc) *API {
	a.Doc = lo.FromPtr(doc)
	return a
}

// SetGuard 设置API的Guard
func (a *API) SetGuard(guard sdwebapp.RouteGuard) *API {
	a.Guard = guard
	return a
}

// SetGuardErrorHandler 设置API的GuardErrorHandler
func (a *API) SetGuardErrorHandler(handler echo.HTTPErrorHandler) *API {
	a.GuardErrorHandler = handler
	return a
}

// AddMiddlewares 添加API的中间件
func (a *API) AddMiddlewares(middlewares ...echo.MiddlewareFunc) *API {
	a.Middlewares = append(a.Middlewares, middlewares...)
	return a
}

// ToRoutes 实现Routable接口
func (a *API) ToRoutes(_ *sdwebapp.App) []*sdwebapp.Route {
	return []*sdwebapp.Route{{
		Name:              a.Name,
		Method:            http.MethodPost,
		Path:              a.Path,
		Handler:           a.Handler,
		Middlewares:       append(a.Middlewares, setRenderResultMiddleware),
		Guard:             a.Guard,
		GuardErrorHandler: a.GuardErrorHandler,
		Doc:               a.Doc,
	}}
}

// Apply 实现Component接口
func (a *API) Apply(app *sdwebapp.App) error {
	return sdwebapp.ApplyRoutes(app, a.ToRoutes(app))
}

func renderResult(c echo.Context, r *sdwebapp.Result) error {
	var r1 Result
	switch d0 := r.Data.(type) {
	case nil:
		r1.Code, r1.Data = CodeOK, nil
	case ResultInterface:
		r1.assignFrom(d0)
	default:
		r1.Code, r1.Data = CodeOK, d0
	}
	switch e0 := r.Err.(type) {
	case nil:
	case ResultInterface:
		r1.assignFrom(e0)
	default:
		if he, ok := sderr.As[*echo.HTTPError](r.Err); ok && he != nil {
			r1.Code = httpStatusCodeToResultCode(he.Code)
			r1.Message = httpErrorMessageToResultMessage(he.Code, he.Message)
		} else {
			r1.Code, r1.Message = CodeUnknown, sderr.PublicMsg(r.Err)
		}
	}
	r1.trimMeta()
	return c.JSON(http.StatusOK, &r1)
}
