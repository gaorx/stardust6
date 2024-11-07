package sdcodegen

import (
	"github.com/gaorx/stardust6/sderr"
	"github.com/samber/lo"
	"io/fs"
	"path/filepath"
	"slices"
)

type (
	// Handler 文本内容生成器
	Handler func(c *Context)

	// Middleware 生成文本内容时的中间件
	Middleware func(c *Context, next Handler)
)

// Generate 生成一个文本文件
func Generate(name string, current *File, h Handler, middlewares ...Middleware) (*File, error) {
	c := Context{current: current, logger: NoLog()}
	middlewares = append([]Middleware{setupStyle}, middlewares...)
	c.name = name
	h1 := func(c1 *Context) {
		h(c1)
		c.expandPlaceholders()
	}
	genErr, ok := lo.TryWithErrorValue(func() error {
		wrap(h1, middlewares)(&c)
		return nil
	})
	if !ok && genErr != nil {
		if _, isAbort := genErr.(abortError); !isAbort {
			return nil, sderr.With("file", name).Wrapf(sderr.Ensure(genErr), "generate file failed")
		}
	}
	if c.err != nil {
		return nil, sderr.With("file", name).Wrapf(c.err, "generate file failed")
	}
	mode := 0600
	if c.executable {
		mode = 0700
	}
	return &File{
		Name:      c.name,
		Data:      c.buff.Bytes(),
		Mode:      fs.FileMode(mode),
		Discarded: c.discarded,
	}, nil
}

// GenerateFile 从指定文件中读取，然后重新生成，最后再写入到这个文件中
func GenerateFile(dirname, name string, h Handler, middlewares ...Middleware) error {
	current, err := ReadFile(dirname, name)
	if err != nil && !IsNotExistErr(err) {
		return sderr.Wrap(err)
	}
	f, err := Generate(name, current, h, middlewares...)
	if err != nil {
		return sderr.Wrap(err)
	}
	err = f.Write(dirname, nil)
	return sderr.Wrap(err)
}

// GenerateBytes 生成一个字节数组形式文件，并返回其内容
func GenerateBytes(h Handler, middlewares ...Middleware) ([]byte, error) {
	tmpf, err := Generate("tmp", nil, h, middlewares...)
	if err != nil {
		return nil, sderr.Wrapf(err, "generate text failed in context")
	}
	return tmpf.Data, nil
}

// GenerateText 生成一个文本文件，并返回其内容
func GenerateText(h Handler, middlewares ...Middleware) (string, error) {
	b, err := GenerateBytes(h, middlewares...)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

type abortError struct{}

func wrap(h Handler, middlewares []Middleware) Handler {
	one := func(h Handler, m Middleware) Handler {
		return func(p *Context) {
			m(p, h)
		}
	}
	inner := func(p *Context) {
		if h != nil {
			h(p)
		}
	}
	h1 := inner
	for i := len(middlewares) - 1; i >= 0; i-- {
		m := middlewares[i]
		if m != nil {
			h1 = one(h1, m)
		}
	}
	return h1
}

func setupStyle(c *Context, next Handler) {
	tabExts := []string{".mk", ".go"}
	space2Exts := []string{
		// html
		".html",
		".htm",
		".haml",
		".tpl",
		".tmpl",
		".gohtml",

		// style
		".css",
		".scss",
		".sass",
		".less",

		// source
		".js",
		".ts",
		".jsx",
		".tsx",
		".coffee",
		".kt",
		".kts",
		".gradle",
		".rs",
		".md",

		// data & config
		".xml",
		".json",
		".yaml",
		".yml",
		".toml",
	}

	c.SetNewline("\n")
	name := c.Name()
	ext := filepath.Ext(name)
	if name == "Makefile" || slices.Contains(tabExts, ext) {
		c.SetTab("\t")
	} else if slices.Contains(space2Exts, ext) {
		c.SetTab("  ")
	} else {
		c.SetTab("    ")
	}
	next(c)
}
