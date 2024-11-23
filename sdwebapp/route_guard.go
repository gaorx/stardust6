package sdwebapp

import (
	"github.com/expr-lang/expr"
	"github.com/gaorx/stardust6/sderr"
	"github.com/gaorx/stardust6/sdsecurity/sdauthn"
	"github.com/labstack/echo/v4"
)

type RouteGuard func(c echo.Context) error

func PermitAll() RouteGuard {
	return func(c echo.Context) error {
		return nil
	}
}

func RejectAll() RouteGuard {
	return func(c echo.Context) error {
		return echo.ErrForbidden
	}
}

func Authenticated() RouteGuard {
	return func(c echo.Context) error {
		err := C(c).Authenticator().AccessAuthenticated()
		if err != nil {
			return sderr.Wrap(err)
		}
		return nil
	}
}

func Authorized(f func(c echo.Context, p *sdauthn.Principal) bool) RouteGuard {
	return func(c echo.Context) error {
		err := C(c).Authenticator().AccessAuthenticated()
		if err != nil {
			return sderr.Wrap(err)
		}
		if f != nil {
			p := C(c).AccessPrincipal()
			ok := f(c, p)
			if !ok {
				return echo.ErrForbidden
			}
		}
		return nil
	}
}

func AuthorizedAuthorities(authorities ...string) RouteGuard {
	return Authorized(func(_ echo.Context, p *sdauthn.Principal) bool {
		return p.HasAnyAuthority(authorities...)
	})
}

func AuthorizedExpr(expression string) RouteGuard {
	type envVars struct {
		ID          string
		Username    string
		Email       string
		Phone       string
		Authorities []string
		Path        string
		Query       string
		QueryParams map[string]string
		PathParams  map[string]string
		Headers     map[string]string
	}

	program, err := expr.Compile(expression, expr.Env(envVars{}))
	if err != nil {
		panic(sderr.With("expr", expression).Wrapf(err, "compile expression failed"))
	}
	return Authorized(func(c echo.Context, p *sdauthn.Principal) bool {
		vars := envVars{
			ID:          p.ID,
			Username:    p.Username,
			Email:       p.Email,
			Phone:       p.Phone,
			Authorities: p.Authorities,
			Path:        c.Path(),
			Query:       c.QueryString(),
			QueryParams: urlValuesToMap(c.QueryParams()),
			PathParams:  C(c).PathParams(),
			Headers:     urlValuesToMap(c.Request().Header),
		}
		r, err := expr.Run(program, vars)
		if err != nil {
			return false
		}
		if b, ok := r.(bool); ok {
			return b
		} else {
			return false
		}
	})
}
