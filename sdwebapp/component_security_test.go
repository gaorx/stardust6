package sdwebapp

import (
	"github.com/gaorx/stardust6/sdjson"
	"github.com/gaorx/stardust6/sdsecurity/sdauthn"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestSecurity(t *testing.T) {
	is := assert.New(t)

	type loginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	loginHandler := func(c Context, req *loginRequest) *Result {
		authnReq := sdauthn.NewUsernameAndPassword(req.Username, req.Password)
		principal, err := c.Authenticator().Grant(authnReq)
		if err != nil {
			he := NewHttpErrorFrom(err, "")
			return OK(sdjson.Object{
				"code":    he.Code,
				"message": he.Message,
			})
		}
		token := sdauthn.NewUserToken(principal.ID, time.Now()).In("test")
		return OK(sdjson.Object{
			"code": http.StatusOK,
			"data": c.Authenticator().EncodeToken(token),
		})
	}

	whoamiHandler := func(c Context) *Result {
		p := c.AccessPrincipal()
		return OK(p.Username).Render(AsText)
	}

	app := New(nil)
	app.MustInstall(
		Security{
			Principals: SimpleUsers([]*SimpleUser{
				{
					ID:       "uid1",
					Username: "user1",
					Password: "123456",
				},
				{
					ID:          "uid2",
					Username:    "user2",
					Password:    "654321",
					Authorities: []string{"ROLE1"},
				},
			}),
			AccessRequest: AccessRequestFromQueryParam("_token", nil),
			RequestCodec:  sdauthn.JWTUserToken("7EH2TaNLUWNA23OE"),
			SignatureVerifiers: SignatureVerifiers{
				VerifySign(func(c echo.Context) bool {
					return true
				}).For("/api/"),
			},
		},
		Routes{
			R("POST", "/api/login", loginHandler),
			R("POST", "/api/whoami", whoamiHandler).SetGuard(HasAuthority("ROLE1")),
			R("POST", "/api/hello2", func(c Context) *Result {
				p := c.AccessPrincipal()
				return OK(p.Username).Render(AsText)
			}).SetGuard(HasAuthority("ROLE2")),
			R("POST", "/api/hello3", func(c Context) *Result {
				p := c.AccessPrincipal()
				if p == nil {
					return OK("anonymous").Render(AsText)
				}
				return OK(p.Username).Render(AsText)
			}).SetGuard(PermitAll()),
		},
	)

	// 模拟正确登陆user2
	resp := NewTestRequest("POST", "/api/login").SetBodyJson(loginRequest{
		Username: "user2",
		Password: "654321",
	}).Call(app)
	token := resp.BodyJsonObject().Get("data").AsString()
	is.True(
		200 == resp.Code &&
			resp.BodyJsonObject().Get("code").AsInt() == 200 &&
			token != "",
	)

	// 模拟用错误的用户名登陆
	resp = NewTestRequest("POST", "/api/login").SetBodyJson(loginRequest{
		Username: "invalid_user",
		Password: "654321",
	}).Call(app)
	is.True(
		200 == resp.Code &&
			resp.BodyJsonObject().Get("code").AsInt() == http.StatusUnauthorized,
	)

	// 模拟登陆时密码错误
	resp = NewTestRequest("POST", "/api/login").SetBodyJson(loginRequest{
		Username: "user2",
		Password: "invalid_password",
	}).Call(app)
	is.True(
		200 == resp.Code &&
			resp.BodyJsonObject().Get("code").AsInt() == http.StatusUnauthorized,
	)

	// 登陆时模拟用户名和密码都错误
	resp = app.TestCall(NewTestRequest("POST", "/api/login").SetBodyJson(loginRequest{
		Username: "invalid_user",
		Password: "invalid_password",
	}))
	is.True(
		200 == resp.Code &&
			resp.BodyJsonObject().Get("code").AsInt() == http.StatusUnauthorized,
	)

	// 判断用户返回的身份信息是否正确，token被放置在_token中
	resp = NewTestRequest("POST", "/api/whoami?_token="+token).SetBodyJson(sdjson.Object{}).Call(app)
	is.True(200 == resp.Code && resp.BodyText() == "user2")

	// 判断用户返回的身份信息是否正确，缺少token
	resp = NewTestRequest("POST", "/api/whoami").SetBodyJson(sdjson.Object{}).Call(app)
	is.True(401 == resp.Code)

	// 判断用户返回的身份信息是否正确，token错误的情形
	resp = NewTestRequest("POST", "/api/whoami?_token=invalid_token").SetBodyJson(sdjson.Object{}).Call(app)
	is.True(401 == resp.Code)

	// 此用户没有ROLE2权限，此时验证用户身份成功，但用户无权限，返回403
	resp = NewTestRequest("POST", "/api/hello2?_token="+token).SetBodyJson(sdjson.Object{}).Call(app)
	is.True(403 == resp.Code)

	// 如果一个route被设置为所有人都可以访问，那么可以提供身份信息，也可以不提供身份信息
	resp = NewTestRequest("POST", "/api/hello3?_token="+token).SetBodyJson(sdjson.Object{}).Call(app)
	is.True(200 == resp.Code && resp.BodyText() == "user2")
	resp = NewTestRequest("POST", "/api/hello3").SetBodyJson(sdjson.Object{}).Call(app)
	is.True(200 == resp.Code && resp.BodyText() == "anonymous")

	// 测试BasicAuth
	app = New(nil)
	app.MustInstall(
		Security{
			Principals: SimpleUsers([]*SimpleUser{
				{ID: "uid1", Username: "user1", Password: "123456"},
			}),
			AccessRequest: AccessRequestFromBasicAuth(nil),
			RequestCodec:  sdauthn.JWTUserToken("7EH2TaNLUWNA23OE"),
		},
		Routes{
			R("POST", "/api/whoami", whoamiHandler),
		},
	)
	resp = NewTestRequest("POST", "/api/whoami").
		SetBasicAuth("user1", "123456").
		SetBodyJson(sdjson.Object{}).
		Call(app)
	is.True(200 == resp.Code && resp.BodyText() == "user1")

	// 测试BearerAuth
	app = New(nil)
	app.MustInstall(
		Security{
			Principals: SimpleUsers([]*SimpleUser{
				{ID: "uid1", Username: "user1", Password: "123456"},
			}),
			AccessRequest: AccessRequestFromBearerAuth(nil),
			RequestCodec:  sdauthn.JWTUserToken("7EH2TaNLUWNA23OE"),
		},
		Routes{
			R("POST", "/api/login", loginHandler),
			R("POST", "/api/whoami", whoamiHandler),
		},
	)
	resp = NewTestRequest("POST", "/api/login").SetBodyJson(loginRequest{
		Username: "user1",
		Password: "123456",
	}).Call(app)
	token = resp.BodyJsonObject().Get("data").AsString()
	is.True(
		200 == resp.Code &&
			resp.BodyJsonObject().Get("code").AsInt() == 200 &&
			token != "",
	)
	resp = NewTestRequest("POST", "/api/whoami").
		SetBearerAuth(token).
		SetBodyJson(sdjson.Object{}).
		Call(app)
	is.True(200 == resp.Code && resp.BodyText() == "user1")

	// 测试guard
	app = New(nil)
	app.MustInstall(
		Security{
			Principals: SimpleUsers([]*SimpleUser{
				{
					ID:          "uid1",
					Username:    "user1",
					Password:    "123456",
					Authorities: []string{"ROLE1"},
				},
				{
					ID:          "uid2",
					Username:    "user2",
					Password:    "123456",
					Authorities: []string{"ROLE2"},
				},
			}),
			AccessRequest: AccessRequestFromBasicAuth(nil),
			RequestCodec:  sdauthn.JWTUserToken("7EH2TaNLUWNA23OE"),
		},
		Routes{
			Text("/api11", "", "api11").SetGuard(PermitAll()),
			Text("/api12", "", "api12").SetGuard(RejectAll()),
			Text("/api21", "", "api21").SetGuard(Authenticated()),
			Text("/api31", "", "api31").SetGuard(HasAuthority("ROLE1")),
			Text("/api41", "", "api41").SetGuard(IsMatched(`'ROLE1' in Authorities`)),
		},
	)

	// api11使用PermitAll，登陆或者不登陆都可以访问
	resp = NewTestRequest("GET", "/api11").Call(app)
	is.True(200 == resp.Code && resp.BodyText() == "api11")
	resp = NewTestRequest("GET", "/api11").SetBasicAuth("user1", "123456").Call(app)
	is.True(200 == resp.Code && resp.BodyText() == "api11")

	// api12使用RejectAll，登陆或者不登陆都不可以访问
	resp = NewTestRequest("GET", "/api12").Call(app)
	is.True(403 == resp.Code)
	resp = NewTestRequest("GET", "/api12").SetBasicAuth("user1", "123456").Call(app)
	is.True(403 == resp.Code)

	// api21使用Authenticated，必须登陆才可以访问
	resp = NewTestRequest("GET", "/api21").Call(app)
	is.True(401 == resp.Code)
	resp = NewTestRequest("GET", "/api21").SetBasicAuth("user1", "123456").Call(app)
	is.True(200 == resp.Code && resp.BodyText() == "api21")

	// api31使用HasAuthority("ROLE1")，必须登陆且有ROLE1权限才可以访问
	resp = NewTestRequest("GET", "/api31").Call(app)
	is.True(401 == resp.Code)
	resp = NewTestRequest("GET", "/api31").SetBasicAuth("user1", "123456").Call(app)
	is.True(200 == resp.Code && resp.BodyText() == "api31")
	resp = NewTestRequest("GET", "/api31").SetBasicAuth("user2", "123456").Call(app)
	is.True(403 == resp.Code)

	// api41使用IsMatched(`'ROLE1' in Authorities`)，必须登陆且有ROLE1权限才可以访问
	resp = NewTestRequest("GET", "/api41").Call(app)
	is.True(401 == resp.Code)
	resp = NewTestRequest("GET", "/api41").SetBasicAuth("user1", "123456").Call(app)
	is.True(200 == resp.Code && resp.BodyText() == "api41")
	resp = NewTestRequest("GET", "/api41").SetBasicAuth("user2", "123456").Call(app)
	is.True(403 == resp.Code)
}
