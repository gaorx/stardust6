package sdwebapp

import (
	"github.com/labstack/echo/v4"
)

type ResultAddon func(echo.Context)

func SetHeader(key, value string) ResultAddon {
	return func(c echo.Context) {
		c.Response().Header().Set(key, value)
	}
}

func SetHeaders(headers map[string]string) ResultAddon {
	return func(c echo.Context) {
		for k, v := range headers {
			c.Response().Header().Set(k, v)
		}
	}
}
