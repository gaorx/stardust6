package sdblueprint

import (
	"github.com/gaorx/stardust6/sdcodegen"
)

type Generator interface {
	Setup(p *ProjectBuilder)
	Generate(p *Project, cg *sdcodegen.Generator)
}
