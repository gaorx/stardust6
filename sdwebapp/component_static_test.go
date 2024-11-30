package sdwebapp

import (
	"github.com/gaorx/stardust6/sdbytes"
	"github.com/gaorx/stardust6/sdcodegen"
	"github.com/gaorx/stardust6/sdfile"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"slices"
	"testing"
	"testing/fstest"
)

func TestStatic(t *testing.T) {
	is := assert.New(t)

	var (
		indexHtml = `<h1>hello</h1>`
		indexCss  = `body { color: red; }`
		indexJs   = `console.log('hello');`
		logoPng   = lo.Must(sdbytes.FromBase64Std(`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`))
	)

	testFiles := func(app *App) {
		resp := NewTestRequest("GET", "/index.html").Call(app)
		is.True(resp.Code == 200 && resp.BodyText() == indexHtml)
		resp = NewTestRequest("GET", "/asserts/css/index.css").Call(app)
		is.True(resp.Code == 200 && resp.BodyText() == indexCss)
		resp = NewTestRequest("GET", "/asserts/js/index.js").Call(app)
		is.True(resp.Code == 200 && resp.BodyText() == indexJs)
		resp = NewTestRequest("GET", "/assets/img/logo.png").Call(app)
		is.True(resp.Code == 200 && slices.Equal(logoPng, resp.BodyBytes()))
	}

	g := sdcodegen.New()
	g.Add("index.html", sdcodegen.Text(indexHtml))
	g.Add("asserts/css/index.css", sdcodegen.Text(indexCss))
	g.Add("asserts/js/index.js", sdcodegen.Text(indexJs))
	g.Add("assets/img/logo.png", sdcodegen.Bytes(logoPng))
	lo.Must0(sdfile.UseTempDir("", "", func(dirname string) {
		publicDirname := filepath.Join(dirname, "public")
		lo.Must0(g.Generate(publicDirname))

		// 测试目录和一些静态文件
		app := New(nil)
		app.MustInstall(
			Dir("/", publicDirname),
			Text("/a.txt", "", "hello"),
			HTML("/a.html", "hello"),
			CSS("/a.css", "hello"),
			Javascript("/a.js", "hello"),
			File("/assets/js/index2.js", filepath.Join(publicDirname, "asserts/js/index.js")),
		)
		testFiles(app)
		resp := NewTestRequest("GET", "/not-exists").Call(app)
		is.True(resp.Code == 404)
		resp = NewTestRequest("GET", "/a.txt").Call(app)
		is.True(resp.Code == 200 && resp.BodyText() == "hello")
		resp = NewTestRequest("GET", "/a.html").Call(app)
		is.True(resp.Code == 200 && resp.BodyText() == "hello")
		resp = NewTestRequest("GET", "/a.css").Call(app)
		is.True(resp.Code == 200 && resp.BodyText() == "hello")
		resp = NewTestRequest("GET", "/a.js").Call(app)
		is.True(resp.Code == 200 && resp.BodyText() == "hello")
		resp = NewTestRequest("GET", "/assets/js/index2.js").Call(app)
		is.True(resp.Code == 200 && resp.BodyText() == indexJs)

		// 测试如果找不到文件就返回index.html的内容
		app = New(nil)
		app.MustInstall(
			Dir("/", publicDirname).SetFallbackToRootIndexPage(),
		)
		testFiles(app)
		resp = NewTestRequest("GET", "/not-exists").Call(app)
		is.True(resp.Code == 200 && resp.BodyText() == indexHtml)
	}))

	// 测试虚拟fsys中读取文件
	fsys := fstest.MapFS{
		"index.html":            {Data: []byte(indexHtml)},
		"asserts/css/index.css": {Data: []byte(indexCss)},
		"asserts/js/index.js":   {Data: []byte(indexJs)},
		"assets/img/logo.png":   {Data: logoPng},
	}
	app := New(nil)
	app.MustInstall(
		DirFS("/", fsys),
		FileFS("/assets/img/logo2.png", "assets/img/logo.png", fsys),
	)
	testFiles(app)
	resp := NewTestRequest("GET", "/assets/img/logo2.png").Call(app)
	is.True(resp.Code == 200 && slices.Equal(logoPng, resp.BodyBytes()))
}
