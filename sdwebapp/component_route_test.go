package sdwebapp

import (
	"github.com/gaorx/stardust6/sderr"
	"github.com/gaorx/stardust6/sdjson"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRoute(t *testing.T) {
	is := assert.New(t)
	app := New(nil)
	app.MustInstall(
		// 测试route和嵌套
		Text("/a1", "", "a1"),
		Routes{
			R("GET", "/a2", func(c Context) string { return "a2" }),
		},
		Routes{
			Text("/a3", "", "a3"),
			Routes{
				Text("/a4", "", "a4"),
			},
		},
	)
	resp := NewTestRequest("GET", "/a1").Call(app)
	is.True(resp.Code == 200 && resp.BodyText() == "a1")
	resp = NewTestRequest("GET", "/a2").Call(app)
	is.True(resp.Code == 200 && resp.BodyText() == "a2")
	resp = NewTestRequest("GET", "/a3").Call(app)
	is.True(resp.Code == 200 && resp.BodyText() == "a3")
	resp = NewTestRequest("GET", "/a4").Call(app)
	is.True(resp.Code == 200 && resp.BodyText() == "a4")

	type bindStruct struct {
		A string `json:"a"`
	}
	app = New(&Options{Debug: true})
	app.MustInstall(
		// 测试handler
		R("GET", "/h1", echo.HandlerFunc(func(c echo.Context) error {
			return c.String(200, "h1")
		})),
		R("GET", "/h2", func(c echo.Context) error {
			return c.String(200, "h2")
		}),
		R("GET", "/h3", func(c echo.Context) {
			_ = c.String(200, "h3")
		}),
		R("GET", "/h4", func(c Context) error {
			return c.String(200, "h4")
		}),
		R("GET", "/h5", func(c Context) {
			_ = c.String(200, "h5")
		}),
		R("GET", "/h6", func(c echo.Context) string {
			return "h6"
		}),
		R("GET", "/h7", func(c Context) string {
			return "h7"
		}),
		R("GET", "/h8", func(c echo.Context) (string, error) {
			return "h8", nil
		}),
		R("GET", "/h9", func(c Context) (string, error) {
			return "h9", nil
		}),
		R("GET", "/h8a", func(c echo.Context) (string, error) {
			return "h8err", sderr.Newf("h8err")
		}),
		R("GET", "/h9a", func(c Context) (string, error) {
			return "h8err", sderr.Newf("h9err")
		}),
		R("GET", "/h10", func(c echo.Context) *Result {
			return OK("h10").Render(AsText)
		}),
		R("GET", "/h10n", func(c echo.Context) Result {
			return *OK("h10n").Render(AsText)
		}),
		R("GET", "/h10a", func(c echo.Context) *Result {
			return Err("h10a").Render(AsText)
		}),
		R("GET", "/h11", func(c echo.Context) *Result {
			return OK("h11").SetStatusCode(488).Render(AsText)
		}),
		R("GET", "/h11a", func(c echo.Context) *Result {
			return Err("h11a").SetStatusCode(488).Render(AsText)
		}),
		R("GET", "/h12", func(c echo.Context) (*Result, error) {
			return OK("h12").Render(AsText), nil
		}),
		R("GET", "/h12a", func(c echo.Context) (*Result, error) {
			return OK("h12a").Render(AsText), sderr.Newf("h12a")
		}),
		R("GET", "/h12b", func(c echo.Context) (*Result, error) {
			return Err("h12b").Render(AsText).SetStatusCode(488), sderr.Newf("h12b")
		}),
		R("GET", "/h12", func(c echo.Context) (*Result, error) {
			return OK("h12").Render(AsText), nil
		}),
		R("POST", "/h13", func(c echo.Context, req bindStruct) string {
			return req.A
		}),
		R("POST", "/h14", func(req *bindStruct) string {
			return req.A
		}),
		R("POST", "/h15", func(c echo.Context, req sdjson.Object) string {
			return req.Get("a").AsString()
		}),
		R("POST", "/h16", func(req sdjson.Value, c echo.Context) string {
			return req.Get("a").AsString()
		}),
	)
	resp = NewTestRequest("GET", "/h1").Call(app)
	is.True(resp.Code == 200 && resp.BodyText() == "h1")
	resp = NewTestRequest("GET", "/h2").Call(app)
	is.True(resp.Code == 200 && resp.BodyText() == "h2")
	resp = NewTestRequest("GET", "/h3").Call(app)
	is.True(resp.Code == 200 && resp.BodyText() == "h3")
	resp = NewTestRequest("GET", "/h4").Call(app)
	is.True(resp.Code == 200 && resp.BodyText() == "h4")
	resp = NewTestRequest("GET", "/h5").Call(app)
	is.True(resp.Code == 200 && resp.BodyText() == "h5")
	resp = NewTestRequest("GET", "/h6").Call(app)
	is.True(resp.Code == 200 && resp.BodyText() == "h6")
	resp = NewTestRequest("GET", "/h7").Call(app)
	is.True(resp.Code == 200 && resp.BodyText() == "h7")
	resp = NewTestRequest("GET", "/h8").Call(app)
	is.True(resp.Code == 200 && resp.BodyText() == "h8")
	resp = NewTestRequest("GET", "/h9").Call(app)
	is.True(resp.Code == 200 && resp.BodyText() == "h9")
	resp = NewTestRequest("GET", "/h8a").Call(app)
	is.True(resp.Code == 500)
	resp = NewTestRequest("GET", "/h9a").Call(app)
	is.True(resp.Code == 500)
	resp = NewTestRequest("GET", "/h10").Call(app)
	is.True(resp.Code == 200 && resp.BodyText() == "h10")
	resp = NewTestRequest("GET", "/h10n").Call(app)
	is.True(resp.Code == 200 && resp.BodyText() == "h10n")
	resp = NewTestRequest("GET", "/h10a").Call(app)
	is.True(resp.Code == 500)
	resp = NewTestRequest("GET", "/h11").Call(app)
	is.True(resp.Code == 488 && resp.BodyText() == "h11")
	resp = NewTestRequest("GET", "/h11a").Call(app)
	is.True(resp.Code == 488)
	resp = NewTestRequest("GET", "/h12").Call(app)
	is.True(resp.Code == 200 && resp.BodyText() == "h12")
	resp = NewTestRequest("GET", "/h12a").Call(app)
	is.True(resp.Code == 500)
	resp = NewTestRequest("GET", "/h12b").Call(app)
	is.True(resp.Code == 488)
	resp = NewTestRequest("POST", "/h13").SetBodyJson(sdjson.Object{"a": "a1"}).Call(app)
	is.True(resp.Code == 200 && resp.BodyText() == "a1")
	resp = NewTestRequest("POST", "/h14").SetBodyJson(sdjson.Object{"a": "a2"}).Call(app)
	is.True(resp.Code == 200 && resp.BodyText() == "a2")
	resp = NewTestRequest("POST", "/h15").SetBodyJson(sdjson.Object{"a": "a3"}).Call(app)
	is.True(resp.Code == 200 && resp.BodyText() == "a3")
	resp = NewTestRequest("POST", "/h16").SetBodyJson(sdjson.Object{"a": "a4"}).Call(app)
	is.True(resp.Code == 200 && resp.BodyText() == "a4")

	// 一些错误的handler
	is.Panics(func() {
		// 不可绑定的类型
		app.MustInstall(R("GET", "/h101", func(c echo.Context, a int) {}))
	})
	is.Panics(func() {
		// 两个body
		app.MustInstall(R("GET", "/h101", func(c echo.Context, req1, req2 *bindStruct) {}))
	})
	is.Panics(func() {
		// 三个返回值
		app.MustInstall(R("GET", "/h101", func(c echo.Context, req1 *bindStruct) (int, int, error) {
			return 0, 0, nil
		}))
	})
}
