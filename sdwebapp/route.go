package sdwebapp

import (
	"github.com/gaorx/stardust6/sderr"
	"github.com/gaorx/stardust6/sdreflect"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"net/http"
	"reflect"
	"slices"
)

type Route struct {
	Name              string
	Method            string
	Path              string
	Handler           any
	Middlewares       []MiddlewareFunc
	Guard             RouteGuard
	GuardErrorHandler echo.HTTPErrorHandler
	Doc
}

type Routable interface {
	ToRoutes(*App) Routes
}

type Routes []*Route
type Routables []Routable

var (
	_ Routable = (*Route)(nil)
	_ Routable = (Routes)(nil)
	_ Routable = (Routables)(nil)
)

func R(method, path string, handler any) *Route {
	return &Route{Method: method, Path: path, Handler: handler}
}

func RG(method, path string, handler any, guard RouteGuard) *Route {
	return R(method, path, handler).SetGuard(guard)
}

func (r *Route) SetDoc(doc *Doc) *Route {
	r.Doc = lo.FromPtr(doc)
	return r
}

func (r *Route) SetGuard(guard RouteGuard) *Route {
	r.Guard = guard
	return r
}

func (r *Route) SetGuardErrorHandler(handler echo.HTTPErrorHandler) *Route {
	r.GuardErrorHandler = handler
	return r
}

func (r *Route) AddMiddlewares(middlewares ...MiddlewareFunc) *Route {
	r.Middlewares = append(r.Middlewares, middlewares...)
	return r
}

func (r *Route) ToRoutes(*App) Routes {
	return []*Route{r}
}

func (rs Routes) ToRoutes(*App) Routes {
	return slices.Clone(rs)
}

func (rs Routables) ToRoutes(app *App) Routes {
	var routes []*Route
	for _, routable := range rs {
		if routable == nil {
			continue
		}
		routes = append(routes, routable.ToRoutes(app)...)
	}
	return routes
}

type route struct {
	name        string
	method      string
	path        string
	handler     echo.HandlerFunc
	middlewares []echo.MiddlewareFunc
}

func (r *Route) build() (*route, error) {
	r1 := lo.FromPtr(r)
	if r1.Method == "" {
		r1.Method = http.MethodGet
	}
	h1, err := r.makeHandler(r1.Handler)
	if err != nil {
		return nil, sderr.Wrap(err)
	}
	h2 := func(c echo.Context) error {
		guard := r1.Guard
		if guard == nil {
			guard = PermitAll()
		}
		err := guard(c)
		if err != nil {
			guardErrorHandler := r1.GuardErrorHandler
			if guardErrorHandler == nil {
				guardErrorHandler = defaultRouteErrorHandler
			}
			guardErrorHandler(err, c)
			return nil
		}
		return h1(c)
	}
	return &route{
		name:        r1.Name,
		method:      r1.Method,
		path:        r1.Path,
		handler:     h2,
		middlewares: dropNil(r1.Middlewares),
	}, nil
}

func (r *Route) makeHandler(h any) (echo.HandlerFunc, error) {
	switch h1 := h.(type) {
	case nil:
		return defaultRouteHandler, nil
	case func(echo.Context) error:
		return h1, nil
	case echo.HandlerFunc:
		return h1, nil
	case func(Context) error:
		return func(c echo.Context) error {
			return h1(C(c))
		}, nil
	case func(echo.Context) *Result:
		return func(c echo.Context) error {
			res := h1(c)
			return res.render(c)
		}, nil
	case func(Context) *Result:
		return func(c echo.Context) error {
			res := h1(C(c))
			return res.render(c)
		}, nil
	case func(echo.Context):
		return func(c echo.Context) error {
			h1(c)
			return nil
		}, nil
	case func(Context):
		return func(c echo.Context) error {
			h1(C(c))
			return nil
		}, nil
	default:
		hv := sdreflect.RootValueOf(h)
		if hv.Kind() != reflect.Func {
			return nil, sderr.Newf("handler must be a function")
		}

		inMaker, outputProc, err := r.deconstruct(hv)
		if err != nil {
			return nil, sderr.Newf("handler must be a function")
		}
		return func(c echo.Context) error {
			var ins []reflect.Value
			for _, inputMaker := range inMaker {
				v, err1 := inputMaker(c)
				if err1 != nil {
					return err1
				}
				ins = append(ins, v)
			}
			outs := hv.Call(ins)
			return outputProc(c, outs)
		}, nil
	}
}

type handlerInputMaker func(echo.Context) (reflect.Value, error)
type handlerOutputProcessor func(echo.Context, []reflect.Value) error

func (r *Route) deconstruct(hv reflect.Value) ([]handlerInputMaker, handlerOutputProcessor, error) {
	ht := hv.Type()

	// input
	var inputMarkers []handlerInputMaker
	numOfContext, numOfBinding := 0, 0
	for _, t := range sdreflect.Ins(ht) {
		if t == sdreflect.T[echo.Context]() {
			numOfContext += 1
			inputMarkers = append(inputMarkers, func(c echo.Context) (reflect.Value, error) {
				return reflect.ValueOf(c), nil
			})
		} else if t == sdreflect.T[Context]() {
			numOfContext += 1
			inputMarkers = append(inputMarkers, func(c echo.Context) (reflect.Value, error) {
				return reflect.ValueOf(C(c)), nil
			})
		} else if sdreflect.IsStruct(t) {
			numOfBinding += 1
			inputMarkers = append(inputMarkers, func(c echo.Context) (reflect.Value, error) {
				p := reflect.New(t).Interface()
				if err := c.Bind(p); err != nil {
					return reflect.Value{}, sderr.Wrapf(err, "bind input error (struct)")
				}
				return reflect.ValueOf(p).Elem(), nil
			})
		} else if sdreflect.IsStructPtr(t) {
			numOfBinding += 1
			inputMarkers = append(inputMarkers, func(c echo.Context) (reflect.Value, error) {
				p := reflect.New(t.Elem()).Interface()
				if err := c.Bind(p); err != nil {
					return reflect.Value{}, sderr.Wrapf(err, "bind input error (structPtr)")
				}
				return reflect.ValueOf(p).Elem(), nil
			})
		} else if sdreflect.IsMap(t, sdreflect.TString, nil) {
			numOfBinding += 1
			inputMarkers = append(inputMarkers, func(c echo.Context) (reflect.Value, error) {
				m := reflect.MakeMap(t).Interface()
				if err := c.Bind(m); err != nil {
					return reflect.Value{}, sderr.Wrapf(err, "bind input error (map)")
				}
				return reflect.ValueOf(m), nil
			})
		} else {
			return nil, nil, sderr.Newf("illegal handler input type")
		}
	}
	if numOfContext > 1 {
		return nil, nil, sderr.Newf("too many context input")
	}
	if numOfBinding > 1 {
		return nil, nil, sderr.Newf("too many binding input")
	}

	// output
	outTypes := sdreflect.Outs(ht)
	var form int
	if len(outTypes) == 0 {
		form = 11
	} else if len(outTypes) == 1 {
		t := outTypes[0]
		if t == sdreflect.TErr {
			form = 21
		} else if t == sdreflect.T[*Result]() {
			form = 23
		} else if t == sdreflect.T[Result]() {
			form = 24
		} else {
			form = 25
		}
	} else if len(outTypes) == 2 {
		t0, t1 := outTypes[0], outTypes[1]
		if t1 != sdreflect.TErr {
			return nil, nil, sderr.Newf("illegal handler output type")
		}
		if t0 == sdreflect.T[*Result]() {
			form = 33
		} else if t0 == sdreflect.T[Result]() {
			form = 34
		} else {
			form = 35
		}
	} else {
		return nil, nil, sderr.Newf("illegal handler output type")
	}
	outputProcessor := func(c echo.Context, outs []reflect.Value) error {
		switch form {
		case 11:
			return nil
		case 21:
			err0 := outs[0].Interface()
			return sderr.Ensure(err0)
		case 22:
			res0 := outs[0].Interface().(*Result)
			if res0 == nil {
				return nil
			}
			return res0.render(c)
		case 24:
			res0 := outs[0].Interface().(Result)
			return res0.render(c)
		case 25:
			res0 := outs[0].Interface()
			return OK(res0).render(c)
		case 33:
			res0 := outs[0].Interface().(*Result)
			err0 := outs[1].Interface()
			var res1 Result
			if res0 != nil {
				res1 = *res0
			}
			if err0 != nil {
				res1.Err = sderr.Ensure(err0)
			}
			return res1.render(c)
		case 34:
			res0 := outs[0].Interface().(Result)
			err0 := outs[1].Interface()
			if err0 != nil {
				res0.Err = sderr.Ensure(err0)
			}
			return res0.render(c)
		case 35:
			data0 := outs[0].Interface()
			err0 := outs[1].Interface()
			res1 := Of(data0, sderr.Ensure(err0))
			return res1.render(c)
		default:
			panic("can't reach here")
		}
	}

	return inputMarkers, outputProcessor, nil
}
