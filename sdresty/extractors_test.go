package sdresty

import (
	"context"
	"github.com/gaorx/stardust6/sdfile"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func TestExtractors(t *testing.T) {
	is := assert.New(t)

	c := New(nil)
	httpmock.ActivateNonDefault(c.GetClient())
	t.Cleanup(httpmock.DeactivateAndReset)

	httpmock.RegisterResponder("GET", "https://xyz.io/text",
		httpmock.NewStringResponder(200, "hello"))
	httpmock.RegisterResponder("GET", "https://xyz.io/json",
		httpmock.NewStringResponder(200, `[{"msg":"hello"}]`))

	resp, err := GET(context.Background(), c, "https://xyz.io/text", ForResponse())
	is.NoError(err)
	is.Equal("hello", string(resp.Body()))

	d, err := GET(context.Background(), c, "https://xyz.io/text", ForBytesBody())
	is.NoError(err)
	is.Equal([]byte("hello"), d)

	text, err := GET(context.Background(), c, "https://xyz.io/text", ForTextBody())
	is.NoError(err)
	is.Equal("hello", text)

	type respStruct struct {
		Msg string `json:"msg"`
	}
	s1, err := GET(context.Background(), c, "https://xyz.io/json", ForJsonBody[[]respStruct]())
	is.NoError(err)
	is.Equal([]respStruct{{Msg: "hello"}}, s1)
	s2, err := GET(context.Background(), c, "https://xyz.io/json", ForJsonBody[[]*respStruct]())
	is.NoError(err)
	is.Equal([]*respStruct{{Msg: "hello"}}, s2)

	jv, err := GET(context.Background(), c, "https://xyz.io/json", ForJsonValue())
	is.NoError(err)
	is.Equal([]any{map[string]any{"msg": "hello"}}, jv.Interface())

	_ = sdfile.UseTempDir("", "", func(dirname string) {
		fn := filepath.Join(dirname, "hello.txt")
		_, err := GET(context.Background(), c, "https://xyz.io/text", SaveToFile(fn, false))
		is.NoError(err)
		is.FileExists(fn)
		is.Equal("hello", sdfile.ReadTextDef(fn, ""))
	})
}
