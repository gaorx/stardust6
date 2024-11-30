package sdwebapp

import (
	"context"
	"github.com/gaorx/stardust6/sderr"
	"github.com/gaorx/stardust6/sdsecurity/sdauthn"
	"github.com/labstack/echo/v4"
	"time"
)

const (
	akSecurityAuthn           = "sdwebapp.security.authn"
	akSecurityAccessRequest   = "sdwebapp.security.access_request"
	akSecurityAccessPrincipal = "sdwebapp.security.access_principal"
)

type Authenticator interface {
	AccessRequest() sdauthn.Request
	AccessPrincipal() *sdauthn.Principal
	AccessAuthenticated() error
	Grant(req sdauthn.Request) (*sdauthn.Principal, error)
	SaveRequest(req sdauthn.Request) error
	EncodeToken(req sdauthn.Request) string
}

var _ Authenticator = (*authenticatorWithContext)(nil)

type authenticatorWithContext struct {
	a *authenticator
	c echo.Context
}

func (a authenticatorWithContext) AccessRequest() sdauthn.Request {
	return a.a.AccessRequest(a.c)
}

func (a authenticatorWithContext) AccessPrincipal() *sdauthn.Principal {
	return a.a.AccessPrincipal(a.c)
}

func (a authenticatorWithContext) AccessAuthenticated() error {
	return a.a.AccessAuthenticated(a.c)
}

func (a authenticatorWithContext) Grant(req sdauthn.Request) (*sdauthn.Principal, error) {
	return a.a.Grant(a.c, req)
}

func (a authenticatorWithContext) SaveRequest(req sdauthn.Request) error {
	return a.a.SaveRequest(a.c, req)
}
func (a authenticatorWithContext) EncodeToken(req sdauthn.Request) string {
	return a.a.EncodeToken(a.c, req)
}

type authenticator struct {
	securityConfig Security
}

func (a *authenticator) AccessRequest(c echo.Context) sdauthn.Request {
	req, err := a.getAccessRequest(c)
	if err != nil {
		return nil
	}
	return req
}

func (a *authenticator) AccessPrincipal(c echo.Context) *sdauthn.Principal {
	p, err := a.loadAccessPrincipal(c)
	if err != nil {
		return nil
	}
	return p
}

func (a *authenticator) AccessAuthenticated(c echo.Context) error {
	p, err := a.loadAccessPrincipal(c)
	if err != nil {
		return sderr.Wrap(err)
	}
	if p == nil {
		return sderr.Wrap(sdauthn.ErrPrincipalNotFound)
	}
	req := Get[loadable[sdauthn.Request]](c, akSecurityAccessRequest).v
	authenticated, err := a.authenticate(c, req, sdauthn.LoaderOf(p))
	if err != nil {
		return sderr.Wrapf(err, "authenticated principal failed")
	}
	c.Set(akSecurityAccessPrincipal, authenticated)
	return nil
}

func (a *authenticator) Grant(c echo.Context, req sdauthn.Request) (*sdauthn.Principal, error) {
	loader := func(ctx context.Context, pid sdauthn.PrincipalId) (*sdauthn.Principal, error) {
		p0, err := a.securityConfig.Principals(c, &a.securityConfig, req.PrincipalId())
		if err != nil {
			return nil, sderr.Wrapf(err, "load principal failed")
		}
		return p0, nil
	}
	p, err := a.authenticate(c, req, sdauthn.LoaderFunc(loader))
	if err != nil {
		return nil, sderr.Wrapf(err, "grant failed")
	}
	return p, nil
}

func (a *authenticator) SaveRequest(c echo.Context, req sdauthn.Request) error {
	save := a.securityConfig.AccessRequestSave
	if save != nil {
		return sderr.Newf("no save function")
	}
	err := save(c, &a.securityConfig, req)
	if err != nil {
		return sderr.Wrapf(err, "save access request failed")
	}
	return nil
}

func (a *authenticator) EncodeToken(c echo.Context, req sdauthn.Request) string {
	codec := a.securityConfig.RequestCodec
	return codec.Encode(req)
}

func (a *authenticator) getAccessRequest(c echo.Context) (sdauthn.Request, error) {
	req1 := GetOrPut(c, akSecurityAccessRequest, func() loadable[sdauthn.Request] {
		req0, err := a.securityConfig.AccessRequest(c, &a.securityConfig)
		if err != nil {
			return loadErr[sdauthn.Request](sderr.Wrapf(err, "load principal request failed"))
		}
		return loadOk(req0)
	})
	return req1.v, req1.err
}

func (a *authenticator) loadAccessPrincipal(c echo.Context) (*sdauthn.Principal, error) {
	req, err := a.getAccessRequest(c)
	if err != nil {
		return nil, err
	}
	p1 := GetOrPut(c, akSecurityAccessPrincipal, func() loadable[*sdauthn.Principal] {
		if req == nil {
			return loadOk[*sdauthn.Principal](nil)
		}
		id := req.PrincipalId()
		p0, err := a.securityConfig.Principals(c, &a.securityConfig, id)
		if err != nil {
			return loadErr[*sdauthn.Principal](
				sderr.With("principal", id.String()).Wrapf(err, "load principal failed"))
		}
		return loadOk(p0)
	})
	return p1.v, p1.err
}

func (a *authenticator) authenticate(c echo.Context, req sdauthn.Request, loader sdauthn.Loader) (*sdauthn.Principal, error) {
	if req == nil {
		return nil, sderr.Newf("no principal request")
	}
	p1, err := sdauthn.Authenticate(
		c.Request().Context(),
		req,
		loader,
		a.securityConfig.AccessRequestExpiration,
		time.Now(),
	)
	if err != nil {
		return nil, sderr.Wrapf(err, "authenticate failed")
	}
	return p1, nil
}
