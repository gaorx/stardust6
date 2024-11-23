package sdwebapp

import (
	"github.com/gaorx/stardust6/sdsecurity/sdauthn"
	"github.com/labstack/echo/v4"
)

type Security struct {
	Principals              PrincipalLoader
	AccessRequest           AccessRequestFunc
	AccessRequestSave       AccessRequestSaveFunc
	AccessRequestExpiration sdauthn.Expiration
	RequestCodec            sdauthn.RequestCodec
	SignatureVerifiers      SignatureVerifiers
}

func (security Security) Apply(app *App) error {
	if security.AccessRequestExpiration == nil {
		security.AccessRequestExpiration = sdauthn.NoExpiration()
	}
	if security.Principals == nil {
		security.Principals = NoPrincipals()
	}

	// 设置Authenticator
	app.Pre(func(next HandlerFunc) HandlerFunc {
		return func(c echo.Context) error {
			c.Set(akSecurityAuthn, &authenticator{securityConfig: security})
			return next(c)
		}
	})

	// 设置签名校验
	app.Pre(security.SignatureVerifiers.ToMiddleware())
	return nil
}
