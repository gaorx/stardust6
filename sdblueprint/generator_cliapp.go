package sdblueprint

import (
	"github.com/gaorx/stardust6/sdcodegen"
	"github.com/gaorx/stardust6/sdcodegen/sdgogen"
	"path/filepath"
)

type CliAppGenerator struct {
	dir string
}

var _ Generator = (*CliAppGenerator)(nil)

func CliApp(dir string) *CliAppGenerator {
	return &CliAppGenerator{dir: dir}
}

func (bg *CliAppGenerator) To(d string) *CliAppGenerator {
	bg.dir = d
	return bg
}

func (bg *CliAppGenerator) ensureDefault() {
	if bg.dir == "" {
		bg.dir = "cmd/app"
	}
}
func (bg *CliAppGenerator) Setup(_ *ProjectBuilder) {
	bg.ensureDefault()
}

func (bg *CliAppGenerator) Generate(p *Project, cg *sdcodegen.Generator) {
	cg.Add(filepath.Join(bg.dir, "main.go"), func(c *sdcodegen.Context) {
		c.DiscardAndAbortIfExists()
		bg.genApp(sdgogen.C(c))
	}, sdgogen.Formatter().AsMiddleware())
}

func (bg *CliAppGenerator) genApp(c *sdgogen.Context) {
	c.Package("main").Newl()
	c.Import([]string{"fmt"}).Newl()
	c.Func("main", nil, nil, func() {
		c.Tab().Line(`fmt.Println("Hello, World!")`)
	}, nil).Newl()
}
