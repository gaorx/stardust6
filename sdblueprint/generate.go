package sdblueprint

import (
	"github.com/gaorx/stardust6/sdcodegen"
	"github.com/gaorx/stardust6/sderr"
)

func (p *Project) NewGenerator() *sdcodegen.Generator {
	cg := sdcodegen.New()
	for _, g := range p.generators {
		g.Generate(p, cg)
	}
	return cg
}

func (p *Project) Try(currents sdcodegen.Files) (sdcodegen.Files, error) {
	cg := p.NewGenerator()
	return cg.TryAll(currents)
}

func (p *Project) MustTry(currents sdcodegen.Files) sdcodegen.Files {
	outs, err := p.Try(currents)
	if err != nil {
		panic(err)
	}
	return outs
}

func (p *Project) Generate() error {
	root := p.Root()
	if root == "" {
		return sderr.Newf("root not set")
	}
	cg := p.NewGenerator()
	return cg.Generate(root)
}
