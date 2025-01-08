package sdblueprint

import (
	"github.com/gaorx/stardust6/sdcodegen"
	"path/filepath"
)

type ReadmeGenerator struct {
	dir      string
	filename string
}

var _ Generator = (*ReadmeGenerator)(nil)

func Readme() *ReadmeGenerator {
	return &ReadmeGenerator{}
}

func (bg *ReadmeGenerator) To(d string) *ReadmeGenerator {
	bg.dir = d
	return bg
}

func (bg *ReadmeGenerator) File(fn string) *ReadmeGenerator {
	bg.filename = fn
	return bg
}

func (bg *ReadmeGenerator) ensureDefault() {
	if bg.filename == "" {
		bg.filename = "README.md"
	}
}

func (bg *ReadmeGenerator) Setup(_ *ProjectBuilder) {
	bg.ensureDefault()
}

func (bg *ReadmeGenerator) Generate(_ *Project, cg *sdcodegen.Generator) {
	cg.Add(filepath.Join(bg.dir, bg.filename), func(c *sdcodegen.Context) {
		c.DiscardAndAbortIfExists()
		c.Line("# README").Newl()
		c.Line("* TODO")
	})
}
