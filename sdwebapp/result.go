package sdwebapp

import (
	"github.com/gaorx/stardust6/sderr"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Result struct {
	Data     any
	View     string
	Err      error
	Addons   []ResultAddon
	Renderer ResultRenderer

	// HTTP
	StatusCode  int
	ContentType string

	// private
	normalized bool
}

func Of(data any, err error) *Result {
	if err != nil {
		return Err(err).SetData(data)
	}
	return OK(data)
}

func OK(data any) *Result {
	return &Result{Data: data}
}

func Err(err any) *Result {
	if err == nil {
		panic(sderr.Newf("err is nil"))
	}
	return &Result{Err: sderr.Ensure(err)}
}

func (r *Result) SetData(data any) *Result {
	r.Data = data
	return r
}

func (r *Result) SetErr(err error) *Result {
	r.Err = err
	return r
}

func (r *Result) SetView(view string) *Result {
	r.View = view
	return r
}

func (r *Result) Also(f ResultAddon) *Result {
	r.Addons = append(r.Addons, f)
	return r
}

func (r *Result) Render(renderer ResultRenderer) *Result {
	r.Renderer = renderer
	return r
}

func (r *Result) SetStatusCode(statusCode int) *Result {
	r.StatusCode = statusCode
	return r
}

func (r *Result) SetContentType(s string) *Result {
	r.ContentType = s
	return r
}

func (r *Result) Normalize(c echo.Context) *Result {
	if r == nil {
		return nil
	}
	r1 := *r
	if !r1.normalized {
		if r1.Err == nil {
			if r1.StatusCode <= 0 {
				r1.StatusCode = http.StatusOK
			}
		} else {
			var he *echo.HTTPError
			if he0, ok := sderr.As[*echo.HTTPError](r1.Err); ok {
				he = he0
			} else {
				he = NewHttpErrorFrom(r1.Err, "")
			}
			r1.Err = he
			if r1.StatusCode <= 0 {
				r1.StatusCode = he.Code
			}
		}
		r1.normalized = true
	}
	return &r1
}

func (r *Result) render(c echo.Context) error {
	if r == nil {
		return c.NoContent(200)
	}
	renderer := r.Renderer
	if renderer == nil {
		renderer = C(c).ResultRenderer()
	}
	if renderer == nil {
		renderer = renderResultDefault
	}
	for _, addon := range r.Addons {
		if addon != nil {
			addon(c)
		}
	}
	return renderer(c, r)
}
