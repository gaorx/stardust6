package sdwebapp

import (
	"github.com/gaorx/stardust6/sderr"
	"github.com/gaorx/stardust6/sdsecurity/sdauthn"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"strings"
)

func defaultRouteHandler(_ echo.Context) error {
	return nil
}

func dropNil[T HandlerFunc | MiddlewareFunc | *sdauthn.SimpleUser](c []T) []T {
	return lo.Filter[T](c, func(v T, _ int) bool { return v != nil })
}

func sanitizeURI(uri string) string {
	// double slash `\\`, `//` or even `\/` is absolute uri for browsers and by redirecting request to that uri
	// we are vulnerable to open redirect attack. so replace all slashes from the beginning with single slash
	if len(uri) > 1 && (uri[0] == '\\' || uri[0] == '/') && (uri[1] == '\\' || uri[1] == '/') {
		uri = "/" + strings.TrimLeft(uri, `/\`)
	}
	return uri
}

func unwrapHttpError(err error) error {
	rootErr := sderr.Root(err)
	if err1, ok := rootErr.(*echo.HTTPError); ok {
		if err1 != nil {
			return err1
		}
		return nil
	}
	return err
}

func urlValuesToMap(vals map[string][]string) map[string]string {
	m := make(map[string]string, len(vals))
	for k, v := range vals {
		if len(v) > 0 {
			m[k] = v[0]
		} else {
			m[k] = ""
		}
	}
	return m
}

type loadable[T any] struct {
	v      T
	err    error
	loaded bool
}

func loadOk[T any](v T) loadable[T] {
	return loadable[T]{v: v, err: nil, loaded: true}
}

func loadErr[T any](err error) loadable[T] {
	return loadable[T]{v: lo.Empty[T](), err: err, loaded: true}
}
