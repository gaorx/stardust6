package sdtemplate

import (
	"bytes"
	"github.com/gaorx/stardust6/sderr"
	"github.com/samber/lo"
	"html/template"
	"io/fs"
	"path"
	"strings"
)

type htmlExecutor struct {
}

var _ Executor = htmlExecutor{}

func (te htmlExecutor) Exec(tmpl string, data any) (string, error) {
	t, err := template.New("").Parse(tmpl)
	if err != nil {
		return "", sderr.Wrapf(err, "parse html template error")
	}
	buff := bytes.NewBufferString("")
	err = t.Execute(buff, data)
	if err != nil {
		return "", sderr.Wrapf(err, "execute html template error")
	}
	return buff.String(), nil
}

func (te htmlExecutor) ExecOr(tmpl string, data any, def string) string {
	r, err := te.Exec(tmpl, data)
	if err != nil {
		return def
	}
	return r
}

func (te htmlExecutor) MustExec(tmpl string, data any) string {
	s, err := te.Exec(tmpl, data)
	if err != nil {
		panic(err)
	}
	return s
}

func HtmlLoad(fsys fs.FS, name string) (*template.Template, error) {
	data, err := fs.ReadFile(fsys, name)
	if err != nil {
		return nil, sderr.With("name", name).Wrapf(err, "read template error")
	}
	t, err := template.New(name).Parse(string(data))
	if err != nil {
		return nil, sderr.With("name", name).Wrapf(err, "parse template error")
	}
	return t, nil
}

type HtmlLoader struct {
	options   HtmlLoaderOptions
	fsys      fs.FS
	templates []*template.Template
}

type HtmlLoaderOptions struct {
	Eager      bool
	Extensions []string
}

func NewHtmlLoader(fsys fs.FS, opts *HtmlLoaderOptions) (*HtmlLoader, error) {
	if fsys == nil {
		return nil, sderr.Newf("nil fsys for load HTML template")
	}
	opts1 := lo.FromPtr(opts)
	if len(opts1.Extensions) <= 0 {
		opts1.Extensions = []string{".gohtml", ".go.html", ".go.tmpl", ".go.tpl"}
	}
	var templates []*template.Template
	if opts.Eager {
		var filenames []string
		_ = fs.WalkDir(fsys, ".", func(filename string, d fs.DirEntry, err error) error {
			if d == nil {
				return nil
			}
			if d.IsDir() {
				return nil
			}
			basename := path.Base(filename)
			matched := lo.ContainsBy(opts1.Extensions, func(ext string) bool {
				return strings.HasSuffix(basename, ext)
			})
			if matched {
				filenames = append(filenames, filename)
			}
			return nil
		})
		for _, filename := range filenames {
			t, err := HtmlLoad(fsys, filename)
			if err != nil {
				return nil, err
			}
			templates = append(templates, t)
		}
	}
	return &HtmlLoader{options: opts1, fsys: fsys, templates: templates}, nil
}

func (loader *HtmlLoader) Load(name string) (*template.Template, error) {
	if loader.options.Eager {
		for _, t := range loader.templates {
			if t.Name() == name {
				return t, nil
			}
		}
		return nil, sderr.With("name", name).Newf("not found template")
	} else {
		return HtmlLoad(loader.fsys, name)
	}
}
