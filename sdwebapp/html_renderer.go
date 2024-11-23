package sdwebapp

import (
	"github.com/gaorx/stardust6/sderr"
	"github.com/gaorx/stardust6/sdtemplate"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"io"
	"io/fs"
)

type HtmlRenderer struct {
	loader *sdtemplate.HtmlLoader
}

var (
	_ echo.Renderer = (*HtmlRenderer)(nil)
	_ echo.Renderer = HtmlRenderer{}
	_ Component     = (*HtmlRenderer)(nil)
	_ Component     = HtmlRenderer{}
)

func MustHtmlRenderer(fsys fs.FS, eager bool) *HtmlRenderer {
	return lo.Must(NewHtmlRenderer(fsys, eager))
}

func NewHtmlRenderer(fsys fs.FS, eager bool) (*HtmlRenderer, error) {
	loader, err := sdtemplate.NewHtmlLoader(fsys, &sdtemplate.HtmlLoaderOptions{
		Eager: eager,
	})
	if err != nil {
		return nil, sderr.Wrap(err)
	}
	return &HtmlRenderer{loader: loader}, nil
}

func (renderer HtmlRenderer) Apply(app *App) error {
	app.Renderer = renderer
	return nil
}

func (renderer HtmlRenderer) Render(wr io.Writer, name string, data any, ec echo.Context) error {
	t, err := renderer.loader.Load(name)
	if err != nil {
		return sderr.Wrap(err)
	}
	err = t.Execute(wr, data)
	if err != nil {
		return sderr.With("name", name).Wrapf(err, "execute template error")
	}
	return nil
}
