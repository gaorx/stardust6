package sdwebapp

import (
	"github.com/gaorx/stardust6/sderr"
	"github.com/gaorx/stardust6/sdsecurity/sdauthn"
)

func (c Context) Authenticator() Authenticator {
	a := Get[*authenticator](c, akSecurityAuthn)
	if a == nil {
		panic(sderr.Newf("authenticator not found"))
	}
	return authenticatorWithContext{a: a, c: c}
}

func (c Context) AccessRequest() sdauthn.Request {
	a := Get[*authenticator](c, akSecurityAuthn)
	if a == nil {
		return nil
	}
	return a.AccessRequest(c)
}

func (c Context) AccessPrincipal() *sdauthn.Principal {
	a := Get[*authenticator](c, akSecurityAuthn)
	if a == nil {
		return nil
	}
	return a.AccessPrincipal(c)
}

func (c Context) AllowCrossOrigin(origin string) {
	c.Response().Header().Add("Access-Control-Allow-Origin", origin)
	c.Response().Header().Add("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	c.Response().Header().Add("Access-Control-Allow-Methods", "PUT,POST,GET,DELETE,OPTIONS")
}
