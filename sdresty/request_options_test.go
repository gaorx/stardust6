package sdresty

import (
	"context"
	"github.com/gaorx/stardust6/sdjson"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRequestOptions(t *testing.T) {
	is := assert.New(t)

	c := New(nil)
	j, err := POST(
		context.Background(),
		c, "https://httpbin.org/{method}",
		ForJsonValue(),
		JsonData(sdjson.Object{"post_k1": "post_v1", "post_k2": 22}),
		QueryParam("k1", 11),
		QueryParams(map[string]any{"k2": 22, "k3": 33}),
		PathParam("method", "post"),
		PathParams(map[string]any{"method": "post"}),
		Header("Custom-Header-1", "xyz1"),
		Headers(map[string]string{"Custom-Header-2": "xyz2"}),
		UserAgent("my-user-agent"),
	)
	is.NoError(err)
	is.Equal("11", j.Get("args", "k1").AsString())
	is.Equal("22", j.Get("args", "k2").AsString())
	is.Equal("33", j.Get("args", "k3").AsString())
	is.Equal("post_v1", j.Get("json", "post_k1").AsString())
	is.Equal(22, j.Get("json", "post_k2").AsInt())
	is.Equal("xyz1", j.Get("headers", "Custom-Header-1").AsString())
	is.Equal("xyz2", j.Get("headers", "Custom-Header-2").AsString())
	is.Equal("application/json", j.Get("headers", "Content-Type").AsString())
	is.Equal("my-user-agent", j.Get("headers", "User-Agent").AsString())

	j, err = POST(
		context.Background(),
		c, "https://httpbin.org/post",
		ForJsonValue(),
		FormData(map[string]any{"form_k1": "post_v1", "form_k2": 22}),
		FileData("file1", "file.txt", []byte("file content")),
	)
	is.NoError(err)
	is.Equal("post_v1", j.Get("form", "form_k1").AsString())
	is.Equal("22", j.Get("form", "form_k2").AsString())
	is.Equal("file content", j.Get("files", "file1").AsString())
}
