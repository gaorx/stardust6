package sdgogen

import (
	"fmt"
	"github.com/gaorx/stardust6/sdcodegen"
	"github.com/samber/lo"
	"sort"
	"strings"
)

type Context struct {
	sdcodegen.ContextWrapper[*Context]
}

func C(c *sdcodegen.Context) *Context {
	c1 := &Context{}
	c1.ContextWrapper = sdcodegen.MakeContextWrapper(c, c1)
	return c1
}

func (c *Context) PrintWarning(times int) *Context {
	for i := 0; i < times; i++ {
		c.Line("// The file is automatically generated. DO NOT EDIT.")
	}
	return c
}

func (c *Context) Comment(text string) *Context {
	c.Print("// ").Line(text)
	return c
}

func (c *Context) Commentf(format string, args ...any) *Context {
	return c.Comment(fmt.Sprintf(format, args...))
}

func (c *Context) CommentNR(text string) *Context {
	c.Print("// ").Print(text)
	return c
}

func (c *Context) CommentfNR(format string, args ...any) *Context {
	return c.CommentNR(fmt.Sprintf(format, args...))
}

func (c *Context) Package(name string) *Context {
	c.Linef("package %s", name)
	return c
}

func (c *Context) Import(pkgs []string) *Context {
	pkgs = lo.Uniq(pkgs)
	sort.Strings(pkgs)
	switch len(pkgs) {
	case 0:
	default:
		c.Line("import (")
		for _, pkg := range pkgs {
			c.Tab().Linef("\"%s\"", pkg)
		}
		c.Line(")")
	}
	return c
}

func (c *Context) Struct(name string, body func()) *Context {
	return c.genStruct(name, body, true)
}

func (c *Context) genStruct(name string, body func(), newline bool) *Context {
	if name != "" {
		c.Linef("type %s struct {", name)
	} else {
		c.Line("struct {")
	}
	if body != nil {
		body()
	}
	c.Print("}")
	if newline {
		c.Newl()
	}
	return c
}

func (c *Context) AnonymousStruct(body func()) *Context {
	return c.genStruct("", body, false)
}

func (c *Context) AnonymousStructPtr(body func()) *Context {
	c.Print("*").AnonymousStruct(body)
	return c
}

func (c *Context) Field(name string, typ string, tags []string, comment string) *Context {
	if name != "" && typ != "" {
		c.Printf("%s %s", name, typ)
	} else if name != "" {
		c.Print(name)
	} else {
		c.Print(typ)
	}
	if len(tags) > 0 {
		joinedTag := strings.Join(tags, " ")
		c.Printf(" `%s`", joinedTag)
	}
	if comment != "" {
		c.Print(" ").CommentNR(comment)
	}
	c.Newl()
	return c
}

type FuncOptions struct {
	MultilineParams  bool
	MultilineReturns bool
}

func (c *Context) Func(name string, params Params, returns Params, body func(), opts *FuncOptions) *Context {
	opts1 := lo.FromPtr(opts)
	c.
		Printf("func %s(%s)", name, params.String(opts1.MultilineParams)).
		Print(returns.StringReturns(opts1.MultilineReturns)).
		Line(" {")
	if body != nil {
		body()
	}
	c.Line("}")
	return c
}

func (c *Context) FuncE(name string, params Params, returnsWithoutErr Params, body func(), opts *FuncOptions) *Context {
	return c.Func(name, params, returnsWithoutErr.WithErr(), body, opts)
}

func (c *Context) MemberFunc(name string, receiver Param, params Params, returns Params, body func(), opts *FuncOptions) *Context {
	opts1 := lo.FromPtr(opts)
	c.
		Printf("func (%s) %s(%s)", receiver.String(), name, params.String(opts1.MultilineParams)).
		Print(returns.StringReturns(opts1.MultilineReturns)).
		Line(" {")
	if body != nil {
		body()
	}
	c.Line("}")
	return c
}

func (c *Context) MemberFuncE(name string, receiver Param, params Params, returnsWithoutErr Params, body func(), opts *FuncOptions) *Context {
	return c.MemberFunc(name, receiver, params, returnsWithoutErr.WithErr(), body, opts)
}
