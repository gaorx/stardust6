package sdwebapp

import (
	"net/http/httptest"
)

func (app *App) TestCall(req *TestRequest) *TestResponse {
	rec := httptest.NewRecorder()
	app.ServeHTTP(rec, req.Request)
	return &TestResponse{rec}
}
