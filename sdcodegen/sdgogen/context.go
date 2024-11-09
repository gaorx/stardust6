package sdgogen

import (
	"github.com/gaorx/stardust6/sdcodegen"
)

type Context struct {
	*sdcodegen.Context
}

func C(c *sdcodegen.Context) *Context {
	return &Context{Context: c}
}

func (c *Context) PrintWarning(times int) *Context {
	for i := 0; i < times; i++ {
		c.Line("// The file is automatically generated. DO NOT EDIT.")
	}
	return c
}

func (c *Context) Package(name string) *Context {
	c.Linef("package %s", name)
	return c
}

func (c *Context) Import(pkgs []string) *Context {
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

func (c *Context) Func(name string, params Params, returns Params, body func()) *Context {
	c.Printf("func %s(%s)", name, params.String()).Print(returns.StringReturns()).Line(" {")
	if body != nil {
		body()
	}
	c.Line("}")
	return c
}

func (c *Context) FuncE(name string, params Params, returnsWithoutErr Params, body func()) *Context {
	return c.Func(name, params, returnsWithoutErr.WithErr(), body)
}

func (c *Context) MemberFunc(name string, receiver Param, params Params, returns Params, body func()) *Context {
	c.Printf("func (%s) %s(%s)", receiver.String(), name, params.String()).Print(returns.StringReturns()).Line(" {")
	if body != nil {
		body()
	}
	c.Line("}")
	return c
}

func (c *Context) MemberFuncE(name string, receiver Param, params Params, returnsWithoutErr Params, body func()) *Context {
	return c.MemberFunc(name, receiver, params, returnsWithoutErr.WithErr(), body)
}
