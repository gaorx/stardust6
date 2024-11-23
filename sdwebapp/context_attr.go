package sdwebapp

import (
	"github.com/gaorx/stardust6/sderr"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
)

func Get[T any](c echo.Context, key string) T {
	var zero T
	return GetOr[T](c, key, zero)
}

func GetOr[T any](c echo.Context, key string, def T) T {
	v0 := c.Get(key)
	if v0 == nil {
		return def
	}
	v1, ok := v0.(T)
	if !ok {
		return def
	}
	return v1
}

func GetOrPut[T comparable](c echo.Context, key string, def func() T) T {
	v := Get[T](c, key)
	if lo.IsNotEmpty(v) {
		return v
	}
	v = def()
	if lo.IsNotEmpty(v) {
		c.Set(key, v)
	}
	return v
}

func GetOrLoad[T comparable](c echo.Context, key string, load func() (T, error)) (T, error) {
	v := Get[T](c, key)
	if lo.IsNotEmpty(v) {
		return v, nil
	}
	v, err := load()
	if err != nil {
		return lo.Empty[T](), sderr.Wrap(err)
	}
	if lo.IsNotEmpty(v) {
		c.Set(key, v)
	}
	return v, nil
}
