package sdsqlgen

import (
	"fmt"
	"github.com/gaorx/stardust6/sdcodegen"
	"github.com/samber/lo"
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
		c.Line("-- The file is automatically generated. DO NOT EDIT.")
	}
	return c
}

func (c *Context) Comment(text string) *Context {
	c.Print("-- ").Line(text)
	return c
}

func (c *Context) Commentf(format string, args ...any) *Context {
	return c.Comment(fmt.Sprintf(format, args...))
}

func (c *Context) CommentNR(text string) *Context {
	c.Print("-- ").Print(text)
	return c
}

func (c *Context) CommentfNR(format string, args ...any) *Context {
	return c.CommentNR(fmt.Sprintf(format, args...))
}

func (c *Context) CreateTable(tableName string, body func(), opts *CreateTableOptions) *Context {
	opts1 := lo.FromPtr(opts)
	c.Print("CREATE").
		If(opts1.Scope != "", " "+opts1.Scope).
		Print(" TABLE").
		If(opts1.IfNotExists, " IF NOT EXISTS").
		Print(" " + tableName).
		Line(" (")
	if body != nil {
		body()
	}
	c.Print(")").
		If(opts1.PostModifier != "", " "+opts1.PostModifier).
		Iff(opts1.Comment != "" && opts1.Dialect != nil,
			" COMMENT %s", MustLiteral(opts1.Comment, opts1.Dialect)).
		Line(";")
	return c
}

func (c *Context) Field(name, typ string, opts *FieldOptions) *Context {
	opts1 := lo.FromPtr(opts)
	c.Print(name).
		Print(" "+typ).
		If(!opts1.Nullable, " NOT NULL").
		If(opts1.AutoIncr, " AUTO_INCREMENT").
		If(opts1.PrimaryKey, " PRIMARY KEY").
		Iff(opts1.Dialect != nil && opts1.Dialect.AllowDefaultValue(typ) && opts1.DefaultValue != nil,
			" DEFAULT %s", MustLiteral(opts1.DefaultValue, opts1.Dialect)).
		If(opts1.Other != "", " "+opts1.Other).
		Iff(opts1.Comment != "" && opts1.Dialect != nil,
			" COMMENT %s", MustLiteral(opts1.Comment, opts1.Dialect)).
		If(opts1.Comma, ",").
		Newl()
	return c
}

func (c *Context) FieldLine(line string, opts *FieldOptions) *Context {
	opts1 := lo.FromPtr(opts)
	c.Print(line).
		If(opts1.Comma, ",").
		Newl()
	return c
}

func (c *Context) FieldLinef(format string, args []any, opts *FieldOptions) *Context {
	return c.FieldLine(fmt.Sprintf(format, args...), opts)
}
