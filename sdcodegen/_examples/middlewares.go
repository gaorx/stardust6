package main

import (
	"fmt"
	"github.com/gaorx/stardust6/sdcodegen"
)

func main() {
	s, err := sdcodegen.GenerateText(func(c *sdcodegen.Context) {
		c.Line("HELLO WORLD")
	}, sdcodegen.AddHeader("HEADER LINE\n"), sdcodegen.AddFooter("FOOTER LINE\n"))
	if err != nil {
		panic(err)
	}
	fmt.Println(s)
	fmt.Println("==============")

	// 也可以自定义中间件
	s, err = sdcodegen.GenerateText(func(c *sdcodegen.Context) {
		c.Line("HELLO WORLD")
	}, func(c *sdcodegen.Context, next sdcodegen.Handler) {
		c.Line("CUSTOM HEADER")
		next(c)
		c.Line("CUSTOM FOOTER")
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(s)
	fmt.Println("===============")

	// 使用generator的时候，以使用全局中间件
	// 同时在(*Generator).Add()是，也可以加入局部中间件
	g := sdcodegen.New()
	g.SetLogger(sdcodegen.Slog(nil))
	g.Use(func(c *sdcodegen.Context, next sdcodegen.Handler) {
		c.Line("GLOBAL MIDDLEWARE HEADER1")
		next(c)
		c.Line("GLOBAL MIDDLEWARE FOOTER1")
	})
	g.Use(func(c *sdcodegen.Context, next sdcodegen.Handler) {
		c.Line("GLOBAL MIDDLEWARE HEADER2")
		next(c)
		c.Line("GLOBAL MIDDLEWARE FOOTER2")
	})
	g.Add("a.txt", func(c *sdcodegen.Context) {
		c.Line("HELLO WORLD")
	}, func(c *sdcodegen.Context, next sdcodegen.Handler) {
		c.Line("LOCAL MIDDLEWARE HEADER3")
		next(c)
		c.Line("LOCAL MIDDLEWARE FOOTER3")
	}, func(c *sdcodegen.Context, next sdcodegen.Handler) {
		c.Line("LOCAL MIDDLEWARE HEADER4")
		next(c)
		c.Line("LOCAL MIDDLEWARE FOOTER4")
	})
	generatedFile, err := g.TryOne("a.txt", nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(generatedFile.StringText())
}
