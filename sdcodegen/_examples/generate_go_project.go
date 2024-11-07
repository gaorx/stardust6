package main

import (
	"fmt"
	"github.com/gaorx/stardust6/sdcodegen"
	"os"
)

func main() {
	g := sdcodegen.New()
	g.SetLogger(sdcodegen.Slog(nil))                      // 设置日志
	autoHeader := sdcodegen.AddHeader("AUTO-GENERATED\n") // 设置自动添加的头部注释
	g.Add("go.mod", func(c *sdcodegen.Context) {
		c.DiscardAndAbortIfExists() // 这行表示如果目标目录中有此文件，则忽略写入，不会覆盖；否则会生成新文件
		c.Line("XXX")
	})
	g.Add("pkg/lib1/lib.go", func(c *sdcodegen.Context) {
		c.Line("package lib1").Newl()
		c.Line(`import "fmt"`).Newl()
		c.Line("func Foo(s string) {")
		c.Tab().Line(`fmt.Println("lib1:", s)`)
		c.Line("}")
	}, autoHeader)
	g.Add("cmd/app1/main.go", func(c *sdcodegen.Context) {
		c.Line("package main").Newl()
		c.Line(`import "pkg/lib1"`).Newl()
		c.Line("func main() {")
		c.Tab().Line(`lib1.Foo("hello")`)
		c.Line("}")
	}, autoHeader)
	g.Add("README.md", sdcodegen.Text("# README\n\n"))
	g.Add("scripts/tool.sh", func(c *sdcodegen.Context) {
		c.SetExecutable(true)
		c.Line("#!/bin/bash")
		c.Line(`echo "hello"`)
	})
	generatedFiles, err := g.TryAll(nil)
	if err != nil {
		panic(err)
	}

	// 输出files中的文件列表，格式类似ls -l命令的输出
	generatedFiles.LL(os.Stdout, false)
	fmt.Println("")

	// 输出files中的所有文件和内容
	fmt.Println(generatedFiles.StringText())

	// 在目录中生成所有文件
	//g.MustGenerate("</path/to/dir>",)
}
