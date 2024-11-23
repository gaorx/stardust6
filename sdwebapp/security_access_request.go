package sdwebapp

import (
	"github.com/gaorx/stardust6/sderr"
	"github.com/gaorx/stardust6/sdsecurity/sdauthn"
	"github.com/labstack/echo/v4"
	"strings"
)

type (
	AccessRequestFunc     func(echo.Context, *Security) (sdauthn.Request, error)
	AccessRequestSaveFunc func(echo.Context, *Security, sdauthn.Request) error
)

type AccessRequestDomainFunc func(echo.Context) string

func AccessRequestDomainFromQueryParam(name string) AccessRequestDomainFunc {
	return func(c echo.Context) string {
		return c.QueryParam(name)
	}
}

func AccessRequestDomainFromQueryPath(name string) AccessRequestDomainFunc {
	return func(c echo.Context) string {
		return c.Param(name)
	}
}

func AccessRequestDomainFromHeader(name string) AccessRequestDomainFunc {
	return func(c echo.Context) string {
		return c.Request().Header.Get(name)
	}
}

func AccessRequestFromBasicAuth(domain AccessRequestDomainFunc) AccessRequestFunc {
	return func(c echo.Context, _ *Security) (sdauthn.Request, error) {
		username, password, ok := c.Request().BasicAuth()
		if !ok {
			return nil, nil
		}
		d := domain.Domain(c)
		return sdauthn.NewUsernameAndPassword(username, password).In(d), nil
	}
}

func AccessRequestFromBearerAuth(domain AccessRequestDomainFunc) AccessRequestFunc {
	return func(c echo.Context, a *Security) (sdauthn.Request, error) {
		s := c.Request().Header.Get("Authorization")
		if s == "" {
			return nil, nil
		}
		if !strings.HasPrefix(s, "Bearer ") {
			return nil, sderr.Newf("invalid authorization header")
		}
		tokenStr := strings.TrimSpace(strings.TrimPrefix(s, "Bearer "))
		if tokenStr == "" {
			return nil, sderr.Newf("no token")
		}
		return parseUserTokenFrom(c, a, tokenStr, domain)
	}
}

func AccessRequestFromQueryParam(name string, domain AccessRequestDomainFunc) AccessRequestFunc {
	return func(c echo.Context, a *Security) (sdauthn.Request, error) {
		tokenStr := c.QueryParam(name)
		if tokenStr == "" {
			return nil, nil
		}
		return parseUserTokenFrom(c, a, tokenStr, domain)
	}
}

func AccessRequestFromQueryHeader(name string, domain AccessRequestDomainFunc) AccessRequestFunc {
	return func(c echo.Context, a *Security) (sdauthn.Request, error) {
		tokenStr := c.Request().Header.Get(name)
		if tokenStr == "" {
			return nil, nil
		}
		return parseUserTokenFrom(c, a, tokenStr, domain)
	}
}

func parseUserTokenFrom(c echo.Context, a *Security, tokenStr string, domain AccessRequestDomainFunc) (sdauthn.Request, error) {
	if a.RequestCodec == nil {
		return nil, sderr.Newf("no request codec")
	}
	req, err := a.RequestCodec.Decode(tokenStr)
	if err != nil {
		return nil, sderr.Wrapf(err, "decode token failed")
	}
	token := req.(sdauthn.UserToken)
	d := domain.Domain(c)
	if d != "" && token.Domain != "" {
		if d != token.Domain {
			return nil, sderr.Newf("invalid token domain")
		}
	} else if d == "" {
		d = token.Domain
	}
	token = token.In(d)
	return token, nil
}

func (d AccessRequestDomainFunc) Domain(c echo.Context) string {
	if d == nil {
		return ""
	}
	return d(c)
}
