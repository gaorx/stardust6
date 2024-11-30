package main

import (
	"github.com/gaorx/stardust6/sdjson"
	"github.com/gaorx/stardust6/sdparse"
	"github.com/gaorx/stardust6/sdsecurity/sdauthn"
	"github.com/gaorx/stardust6/sdslog"
	"github.com/gaorx/stardust6/sdwebapp"
	"github.com/gaorx/stardust6/sdwebapp/sdsimpleapi"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"log/slog"
	"testing/fstest"
	"time"
)

func main() {
	// 设置默认的slog logger
	sdslog.SetDefault([]slog.Handler{
		sdslog.TextFile(slog.LevelDebug, "stdout", true),
	}, nil)

	type N struct {
		Message string
	}

	type User struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		Age  string `json:"age"`
		N    string `json:"n"`
	}

	app := sdwebapp.New(&sdwebapp.Options{
		Debug:         true,
		DisplayBanner: true,
		DisplayPort:   true,
		SlogOptions: &sdwebapp.SlogOptions{
			Logger: slog.Default(), // 使用默认slog的logger进行输出
		},
	})

	// 设置一个mock的目录
	staticFsys := fstest.MapFS{
		"index.html": &fstest.MapFile{Data: []byte("<h1>INDEX.HTML</h1>")},
		"a1.txt":     &fstest.MapFile{Data: []byte("A1")},
		"a2.txt":     &fstest.MapFile{Data: []byte("A2")},
	}

	// 安装各种功能
	app.MustInstall(
		// 向context中注入各种信息
		sdwebapp.Inject(
			// 注入状态
			sdwebapp.State(&N{"xxx"}),
			// 使用Attr注入
			sdwebapp.Attr("k1", "v11"),
			sdwebapp.Attr("k2", "v22"),
			// 使用Attrs注入
			sdwebapp.Attrs{
				sdwebapp.Attr("k3", "v33"),
				sdwebapp.Attr("k4", "v44"),
			},
			// 以map形式注入
			sdwebapp.AttrMap{
				"k5": "k55",
				"k6": "k66",
			},
		),

		// 设置安全信息
		sdwebapp.Security{
			// 设置有两个用户，user1和user2
			Principals: sdwebapp.SimpleUsers([]*sdwebapp.SimpleUser{
				{
					Username:    "user1",
					Password:    "123456",
					Expiry:      sdparse.TimeOr("2024-11-25", time.Time{}),
					Authorities: []string{"ROLE2", "ROLE1"},
				},
				{Username: "user2", Password: "654321"},
			}),
			// Token从query参数的_token中获取
			AccessRequest: sdwebapp.AccessRequestFromQueryParam("_token", nil),
			// 使用JWTUserToken进行请求解码
			RequestCodec: sdauthn.JWTUserToken("7EH2TaNLUWNA23OE"),
			// 设置签名验证器
			SignatureVerifiers: sdwebapp.SignatureVerifiers{
				sdwebapp.VerifySign(func(c echo.Context) bool {
					return true
				}).For("/api/"),
			},
		},

		// 设置一个目录到/static
		sdwebapp.DirFS("/static", staticFsys).SetFallbackToRootIndexPage(),

		// 各种静态route
		sdwebapp.Routes{
			sdwebapp.FileFS("/b1.txt", "/a1.txt", staticFsys),
			sdwebapp.HTML("/aa.html", "<h1>AA.HTML</h1>"),
			sdwebapp.R("GET", "/bb", func(c echo.Context) error {
				return c.String(200, "BB")
			}),
		},

		// 普通的handler
		sdwebapp.R("GET", "/cc", func(c echo.Context) error {
			v1 := sdwebapp.Get[string](c, "k11")
			v2 := sdwebapp.Get[string](c, "k2")
			return c.String(200, "CC/"+v1+"/"+v2)
		}),

		// login的handler
		sdwebapp.R("POST", "/api/login", func(c echo.Context) error {
			var req sdauthn.UsernameAndPassword
			err := c.Bind(&req)
			if err != nil {
				return c.JSON(400, sdjson.Object{"error": "bad request"})
			}
			a := sdwebapp.C(c).Authenticator()
			p, err := a.Grant(&req)
			if err != nil {
				return c.JSON(401, sdjson.Object{"error": err.Error()})
			}
			token := sdauthn.NewUserToken(p.ID, time.Now()).In("app1")
			return c.JSON(200, sdjson.Object{"token": a.EncodeToken(token)})
		}),

		// 一个简单的返回文本的API
		sdwebapp.R("POST", "/api/foo", func(c echo.Context) *sdwebapp.Result {
			return sdwebapp.OK(sdjson.Object{"msg": "foo"}).Render(sdwebapp.AsText)
		}).SetGuard(sdwebapp.IsMatched(`'ROLE1' in Authorities`)),

		// 一个简单的返回User的API
		sdsimpleapi.R("/api/bar", func(c sdwebapp.Context, req struct {
			Name string `json:"name"`
			ID   string `json:"id"`
			Age  string `query:"age"`
		}) (*User, error) {
			return &User{
				ID:   "~" + req.ID,
				Name: "~" + req.Name,
				Age:  req.Age,
				N:    c.State().(*N).Message,
			}, nil
		}, sdwebapp.PermitAll()),
	)

	// 运行
	lo.Must0(app.Start(":8080"))
}
