package sdwebapp

import (
	"github.com/gaorx/stardust6/sderr"
	"github.com/gaorx/stardust6/sdsecurity/sdauthn"
	"github.com/labstack/echo/v4"
	"net/http"
)

func NewHttpErrorFrom(err error, message string) *echo.HTTPError {
	if err == nil {
		return nil
	}
	newHttpErr := func(code int, internal error) *echo.HTTPError {
		if message != "" {
			return echo.NewHTTPError(code, message, internal)
		} else {
			return echo.NewHTTPError(code, httpStatusToText(code))
		}
	}
	if he, ok := sderr.As[*echo.HTTPError](err); ok {
		he1 := *he
		if message != "" {
			he1.Message = message
		} else {
			if he1.Message == nil || he1.Message == "" {
				he1.Message = httpStatusToText(he1.Code)
			}
		}
		return &he1
	} else if sderr.Is(err, sdauthn.ErrPrincipalNotFound) {
		return newHttpErr(http.StatusUnauthorized, err)
	} else if sderr.Is(err, sdauthn.ErrCredentialInvalid) {
		return newHttpErr(http.StatusUnauthorized, err)
	} else if sderr.Is(err, sdauthn.ErrPrincipalDisabled) {
		return newHttpErr(http.StatusUnauthorized, err)
	} else if sderr.Is(err, sdauthn.ErrPrincipalExpired) {
		return newHttpErr(http.StatusUnauthorized, err)
	} else if sderr.Is(err, sdauthn.ErrPrincipalLocked) {
		return newHttpErr(http.StatusUnauthorized, err)
	} else if sderr.Is(err, sdauthn.ErrCredentialExpired) {
		return newHttpErr(http.StatusUnauthorized, err)
	} else {
		return newHttpErr(http.StatusInternalServerError, err)
	}
}

func defaultRouteErrorHandler(err error, c echo.Context) {
	he := NewHttpErrorFrom(err, "")
	if he == nil {
		panic("invalid error")
	}
	_ = c.String(he.Code, sderr.Ensure(he.Message).Error())
}

func httpStatusToText(statusCode int) string {
	msg := http.StatusText(statusCode)
	if msg == "" {
		msg = "Unknown status code"
	}
	return msg
}
