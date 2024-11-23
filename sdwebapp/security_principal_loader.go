package sdwebapp

import (
	"github.com/gaorx/stardust6/sdsecurity/sdauthn"
	"github.com/labstack/echo/v4"
)

type SimpleUser = sdauthn.SimpleUser
type PrincipalLoader func(c echo.Context, a *Security, id sdauthn.PrincipalId) (*sdauthn.Principal, error)

func NoPrincipals() PrincipalLoader {
	return func(_ echo.Context, _ *Security, _ sdauthn.PrincipalId) (*sdauthn.Principal, error) {
		return nil, sdauthn.ErrPrincipalNotFound
	}
}

func SimpleUsers(users []*sdauthn.SimpleUser) PrincipalLoader {
	users = dropNil(users)
	if len(users) <= 0 {
		return NoPrincipals()
	}
	loader := sdauthn.SimpleUsers(users)
	return func(c echo.Context, _ *Security, id sdauthn.PrincipalId) (*sdauthn.Principal, error) {
		return loader.LoadPrincipal(c.Request().Context(), id)
	}
}
