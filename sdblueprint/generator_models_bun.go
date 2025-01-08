package sdblueprint

import (
	"fmt"
	"github.com/gaorx/stardust6/sdcodegen"
	"github.com/gaorx/stardust6/sdcodegen/sdgogen"
	"path/filepath"
	"slices"
	"strings"
)

type ModelsBunGenerator struct {
	dir      string
	filename string
	pred     SchemaPredicate
	pkg      string
}

func ModelsBun() *ModelsBunGenerator {
	return &ModelsBunGenerator{}
}

var _ Generator = (*ModelsBunGenerator)(nil)

func (bg *ModelsBunGenerator) To(d string) *ModelsBunGenerator {
	bg.dir = d
	return bg
}

func (bg *ModelsBunGenerator) File(fn string) *ModelsBunGenerator {
	bg.filename = fn
	return bg
}

func (bg *ModelsBunGenerator) Predicate(pred SchemaPredicate) *ModelsBunGenerator {
	bg.pred = pred
	return bg
}

const (
	goModelsDir      = "pkg/models"
	goModelsFilename = "models.gen.go"
	goModelsPkg      = "models"
)

func (bg *ModelsBunGenerator) ensureDefault() {
	if bg.dir == "" {
		bg.dir = goModelsDir
	}
	if bg.filename == "" {
		bg.filename = goModelsFilename
	}
	if bg.pred == nil {
		bg.pred = ByAll()
	}
	if bg.pkg == "" {
		bg.pkg = lastStringOf(strings.Split(bg.dir, "/"))
	}
	if bg.pkg == "" {
		bg.pkg = goModelsPkg
	}
}

func (bg *ModelsBunGenerator) Setup(b *ProjectBuilder) {
	bg.ensureDefault()
	b.SchemaRefs(bg.pred, func(p *Project, schema Schema) map[string]string {
		goName := schema.Names().Go()
		goRef := fmt.Sprintf("*%s%s.%s", p.modName, ensurePrefix(bg.dir, "/"), goName)
		return map[string]string{"go": goRef}
	})
}

func (bg *ModelsBunGenerator) Generate(p *Project, cg *sdcodegen.Generator) {
	h := func(c *sdcodegen.Context) {
		bg.genTables(sdgogen.C(c), p)
	}
	cg.Add(filepath.Join(bg.dir, bg.filename), h, sdgogen.Formatter().AsMiddleware())
}

func (bg *ModelsBunGenerator) genTables(c *sdgogen.Context, p *Project) {
	tables := p.Tables(bg.pred)
	c.Package(bg.pkg).Newl()
	c.PrintWarning(3).Newl()
	var importPkgs []string
	if len(tables) > 0 {
		importPkgs = append(importPkgs, "github.com/uptrace/bun")
	}
	c.Placeholder("imports").Newl()
	c.ExpandPlaceholder("imports", func() {
		c.Import(importPkgs)
	})
	for _, t := range tables {
		modelImportPkgs := bg.genTable(c, p, t)
		importPkgs = append(importPkgs, modelImportPkgs...)
	}
}

func (bg *ModelsBunGenerator) genTable(c *sdgogen.Context, _ *Project, t *Table) []string {
	var importPkgs []string

	isPk := func(col *Property, indexes []*Index) bool {
		for _, idx := range indexes {
			if idx.Kind() == IndexPK {
				if slices.Contains(idx.Columns(), col.Id()) {
					return true
				}
			}
		}
		return false
	}

	modelName := t.Names().Go()
	c.Commentf("%s Tables %s", modelName, t.Id())
	c.Struct(modelName, func() {
		c.Tab().Field("", "bun.BaseModel", []string{
			fmt.Sprintf(`bun:"table:%s"`, t.Names().Sql()),
		}, "")
		for _, col := range t.Columns() {
			colName := col.Names().Go()
			colRef := goGenRef(c, col.Type())
			importPkgs = append(importPkgs, colRef.Pkgs...)
			fieldTags := []string{
				goMakeJsonTag(col),
				bg.makeBunTag(col, isPk(col, t.Indexes())),
			}
			fieldTags = append(fieldTags, col.Anns().Go()...)
			c.Tab().Field(colName, colRef.Code, fieldTags, col.Comment())
		}
	}).Newl()
	return importPkgs
}

func (bg *ModelsBunGenerator) makeBunTag(col *Property, pk bool) string {
	segs := []string{col.Names().Sql()}
	if pk {
		segs = append(segs, "pk")
	}
	if col.AutoIncr() {
		segs = append(segs, "autoincrement")
	}
	return fmt.Sprintf(`bun:"%s"`, strings.Join(segs, ","))
}
