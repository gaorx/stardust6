package sdsimpleapi

import (
	"github.com/gaorx/stardust6/sderr"
	"github.com/gaorx/stardust6/sdwebapp"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"net/http"
)

type API struct {
	Name              string
	Path              string
	Handler           any
	Middlewares       []echo.MiddlewareFunc
	Guard             sdwebapp.RouteGuard
	GuardErrorHandler echo.HTTPErrorHandler
	sdwebapp.Doc
}

func R(path string, handler any, guard sdwebapp.RouteGuard) *API {
	return &API{Path: path, Handler: handler, Guard: guard}
}

func (a *API) SetName(name string) *API {
	a.Name = name
	return a
}

func (a *API) SetDoc(doc *sdwebapp.Doc) *API {
	a.Doc = lo.FromPtr(doc)
	return a
}

func (a *API) SetGuard(guard sdwebapp.RouteGuard) *API {
	a.Guard = guard
	return a
}

func (a *API) SetGuardErrorHandler(handler echo.HTTPErrorHandler) *API {
	a.GuardErrorHandler = handler
	return a
}

func (a *API) AddMiddlewares(middlewares ...echo.MiddlewareFunc) *API {
	a.Middlewares = append(a.Middlewares, middlewares...)
	return a
}

func (a *API) ToRoutes(*sdwebapp.App) sdwebapp.Routes {
	return sdwebapp.Routes{
		&sdwebapp.Route{
			Name:              a.Name,
			Method:            http.MethodPost,
			Path:              a.Path,
			Handler:           a.Handler,
			Middlewares:       append(a.Middlewares, Middleware),
			Guard:             a.Guard,
			GuardErrorHandler: a.GuardErrorHandler,
			Doc:               a.Doc,
		},
	}
}

func (a *API) Apply(app *sdwebapp.App) error {
	return a.ToRoutes(app).Apply(app)
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
		r1.Code, r1.Message = CodeUnknown, sderr.PublicMsg(r.Err)
	}
	r1.trimMeta()
	return c.JSON(http.StatusOK, &r1)
}
