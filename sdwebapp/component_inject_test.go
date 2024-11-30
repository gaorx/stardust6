package sdwebapp

import (
	"github.com/gaorx/stardust6/sdjson"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInject(t *testing.T) {
	is := assert.New(t)
	app := New(nil)
	app.MustInstall(
		Inject(
			Attr("k1", "v1"),
			Attrs{Attr("k2", "v2")},
			AttrMap{"k3": "v3"},
		),
		R("GET", "/abc", func(c Context) *Result {
			return OK(sdjson.Object{
				"k1": c.Get("k1"),
				"k2": c.Get("k2"),
				"k3": c.Get("k3"),
			}).Render(AsJson)
		}),
	)
	res := NewTestRequest("GET", "/abc").Call(app)
	is.Equal(200, res.Code)
	is.Equal("v1", res.BodyJsonObject().Get("k1").AsString())
	is.Equal("v2", res.BodyJsonObject().Get("k2").AsString())
	is.Equal("v3", res.BodyJsonObject().Get("k3").AsString())
}
