package sdwebapp

import (
	"github.com/gaorx/stardust6/sderr"
	"github.com/gaorx/stardust6/sdsecurity/sdauthn"
	"github.com/labstack/echo/v4"
	"net/http"
)

func NewHttpErrorFrom(err error, retainMessage bool) *echo.HTTPError {
	if err == nil {
		return nil
	}
	newHttpErr := func(code int, internal error) *echo.HTTPError {
		if retainMessage {
			return echo.NewHTTPError(code, internal.Error(), internal)
		} else {
			return echo.NewHTTPError(code, http.StatusText(code), internal)
		}
	}
	if he, ok := sderr.As[*echo.HTTPError](err); ok {
		he1 := *he
		return &he1
	} else if sderr.Is(err, sdauthn.ErrPrincipalNotFound) {
		return newHttpErr(http.StatusUnauthorized, err)
	} else if sderr.Is(err, sdauthn.ErrCredentialInvalid) {
		return newHttpErr(http.StatusBadRequest, err)
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
	he := NewHttpErrorFrom(err, c.Echo().Debug)
	if he == nil {
		panic("invalid error")
	}
	_ = c.String(he.Code, sderr.Ensure(he.Message).Error())
}
