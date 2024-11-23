package sdsimpleapi

import (
	"github.com/gaorx/stardust6/sdwebapp"
	"github.com/labstack/echo/v4"
)

func Middleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sdwebapp.C(c).SetResultRenderer(renderResult)
		return next(c)
	}
}
