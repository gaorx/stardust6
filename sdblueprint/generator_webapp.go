package sdblueprint

import (
	"fmt"
	"github.com/gaorx/stardust6/sdcodegen"
	"github.com/gaorx/stardust6/sdcodegen/sdgogen"
	"github.com/samber/lo"
	"path/filepath"
	"strings"
)

type WebAppGenerator struct {
	dir           string
	withRoutes    bool
	withPrototype bool
}

var _ Generator = (*WebAppGenerator)(nil)

func WebApp(dir string) *WebAppGenerator {
	return &WebAppGenerator{dir: dir}
}

func (bg *WebAppGenerator) To(d string) *WebAppGenerator {
	bg.dir = d
	return bg
}

func (bg *WebAppGenerator) WithRoutes() *WebAppGenerator {
	bg.withRoutes = true
	return bg
}

func (bg *WebAppGenerator) WithPrototype() *WebAppGenerator {
	bg.withPrototype = true
	return bg
}

func (bg *WebAppGenerator) ensureDefault() {
	if bg.dir == "" {
		bg.dir = "cmd/webapp"
	}
}

func (bg *WebAppGenerator) Setup(_ *ProjectBuilder) {
	bg.ensureDefault()
}

func (bg *WebAppGenerator) Generate(p *Project, cg *sdcodegen.Generator) {
	cg.Add(filepath.Join(bg.dir, "main.go"), func(c *sdcodegen.Context) {
		c.DiscardAndAbortIfExists()
		bg.genApp(sdgogen.C(c), p)
	}, sdgogen.Formatter().AsMiddleware())
	if bg.withRoutes {
		cg.Add(filepath.Join(bg.dir, "api_routes.gen.go"), func(c *sdcodegen.Context) {
			bg.genApiRoutes(sdgogen.C(c), p)
		}, sdgogen.Formatter().AsMiddleware())
	}
	if bg.withPrototype {
		for _, api := range p.apis {
			if api.Protocol() == APISimple {
				cg.Add(filepath.Join(bg.dir, bg.getPrototypeFilename(api)), func(c *sdcodegen.Context) {
					c.DiscardAndAbortIfExists()
					bg.genApiPrototype(sdgogen.C(c), api)
				}, sdgogen.Formatter().AsMiddleware())
			}
		}
	}
}

func (bg *WebAppGenerator) genApp(c *sdgogen.Context, _ *Project) {
	c.Package("main").Newl()
	c.Import([]string{
		"github.com/gaorx/stardust6/sdslog",
		"github.com/gaorx/stardust6/sdwebapp",
		"log/slog",
	}).Newl()
	c.Func("main", nil, nil, func() {
		c.Tab().Line(`sdslog.SetDefault([]sdslog.Handler{`)
		c.TabX(2).Line(`sdslog.TextFile(slog.LevelDebug, "stdout", true),`)
		c.Tab().Line(`}, nil)`)
		c.Tab().Line(`app := sdwebapp.New(&sdwebapp.Options{`)
		c.TabX(2).Line(`Debug: true,`)
		c.Tab().Line(`})`)
		c.Tab().Line(`app.MustInstall(`)
		c.TabX(2).Line(`sdwebapp.R("GET", "/ping", func(c sdwebapp.Context) *sdwebapp.Result {`)
		c.TabX(3).Line(`return sdwebapp.OK("pong")`)
		c.TabX(2).Line(`}),`)
		c.If(bg.withRoutes, func() {
			c.TabX(2).Line(`apiRoutes,`)
		})
		c.Tab().Line(`)`)
		c.Tab().Line(`err := app.Start(":8080")`)
		c.Tab().Line(`if err != nil {`)
		c.TabX(2).Line(`slog.With(sdslog.E(err)).Error("start app failed")`)
		c.Tab().Line(`}`)
	}, nil).Newl()
}

func (bg *WebAppGenerator) genApiRoutes(c *sdgogen.Context, p *Project) {
	c.Package("main").Newl()
	c.PrintWarning(3).Newl()
	c.Import([]string{
		"github.com/gaorx/stardust6/sdwebapp",
		"github.com/gaorx/stardust6/sdwebapp/sdsimpleapi",
	}).Newl()
	c.Line(`apiRoutes := sdwebapp.Routes{`)
	p.APIs().ForEach(func(api *APIDecl) {
		switch api.Protocol() {
		case APISimple:
			c.Tab().Linef(`sdsimpleapi.R("%s", %s, %s),`, api.HttpPath(), bg.getRouteFuncName(api), bg.getGuardCode(api))
		}
	})
	c.Line(`}`)
}

func (bg *WebAppGenerator) genApiPrototype(c *sdgogen.Context, api *APIDecl) {
	var importPkgs []string
	c.Package("main").Newl()
	c.Placeholder("imports")
	c.ExpandPlaceholder("imports", func() {
		c.Import(importPkgs)
	})
	inRef, outRef := goGenRef(c, api.In()), goGenRef(c, api.Out())
	importPkgs = append(importPkgs, inRef.Pkgs...)
	importPkgs = append(importPkgs, outRef.Pkgs...)
	c.FuncE(
		bg.getRouteFuncName(api),
		sdgogen.Params{
			sdgogen.P("c", "sdwebapp.Context"),
			sdgogen.P("input", inRef.Code),
		},
		sdgogen.Params{
			sdgogen.P("", outRef.Code),
		},
		func() {
			c.Tab().Line(`panic("TODO")`)
		},
		&sdgogen.FuncOptions{
			MultilineParams:  strings.Contains(inRef.Code, "struct"),
			MultilineReturns: strings.Contains(outRef.Code, "struct"),
		},
	)
}

func (bg *WebAppGenerator) getRouteFuncName(api *APIDecl) string {
	s := ToCamel(api.Id())
	s = strings.TrimPrefix(s, "api")
	s = strings.TrimSuffix(s, "API")
	s = strings.TrimSuffix(s, "Api")
	s = s + "API"
	return s
}

func (bg *WebAppGenerator) getGuardCode(api *APIDecl) string {
	switch api.Guard() {
	case APIPermitAll:
		return "sdwebapp.PermitAll()"
	case APIRejectAll:
		return "sdwebapp.RejectAll()"
	case APIAuthenticated:
		return "sdwebapp.Authenticated()"
	case APIHasAuthority:
		authorities := strings.Join(lo.Map(api.GuardVals(), func(authority string, _ int) string {
			return "\"" + authority + "\""
		}), ", ")
		return "sdwebapp.HasAuthority(" + authorities + ")"
	case APIIsMatched:
		return "sdwebapp.IsMatched(\"" + api.GuardVal() + "\")"
	default:
		panic("unknown guard")
	}
}

func (bg *WebAppGenerator) getPrototypeFilename(api *APIDecl) string {
	pathSegs := lo.Filter(strings.Split(api.HttpPath(), "/"), func(s string, _ int) bool {
		return s != "" && s != "/"
	})
	return fmt.Sprintf("z_%s.go", strings.Join(pathSegs, "_"))
}
