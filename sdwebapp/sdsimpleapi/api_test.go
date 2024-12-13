package sdsimpleapi

import (
	"github.com/gaorx/stardust6/sdwebapp"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAPI(t *testing.T) {
	is := assert.New(t)

	type request struct {
		Name string `json:"name"`
	}

	app := sdwebapp.New(nil)
	app.MustInstall(
		R("/api/foo", func(req request) *Result {
			return OK("hello, " + req.Name)
		}, sdwebapp.PermitAll()),
		R("/api/bar", func(req request) *Result {
			return Err("SOME_ERROR", "xxx")
		}, sdwebapp.PermitAll()),
	)

	// foo
	NewTestRequest("/api/foo", request{
		Name: "world",
	}).Call(app).Let(func(resp *sdwebapp.TestResponse) {
		var r ResultT[string]
		is.Equal(200, resp.Code)
		is.NotPanics(func() {
			resp.BodyJson(&r)
		})
		is.Equal(CodeOK, r.Code)
		is.Equal("hello, world", r.Data)
	})

	// bar
	NewTestRequest("/api/bar", request{
		Name: "world",
	}).Call(app).Let(func(resp *sdwebapp.TestResponse) {
		var r ResultT[string]
		is.Equal(200, resp.Code)
		is.NotPanics(func() {
			resp.BodyJson(&r)
		})
		is.Equal("SOME_ERROR", r.Code)
		is.Equal("xxx", r.Message)
	})
}
