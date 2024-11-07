package main

import (
	"github.com/gaorx/stardust6/sdcodegen"
	"time"
)

func main() {
	g := sdcodegen.New()
	g.SetLogger(sdcodegen.Slog(nil))
	g.Add("a.txt", func(c *sdcodegen.Context) {
		// 因为在generate时并不读取此文件，所以其中的current总是nil
		c.Line("HELLO ", time.Now().Format(time.DateTime))
	})
	g.Add("b.sh", func(c *sdcodegen.Context) {
		// 因为在generate时读取此文件，所以其中的current在第一次生成时是nil，后续是当前的文件内容
		c.SetExecutable(true)
		c.Line("#!/bin/bash")
		c.Linef("echo 'hello %s'", time.Now().Format(time.DateTime))
	})
	// 下面使用了sdcodegen.HasSuffix用来指定只加载当前存在的.sh和.go文件，忽略其他文件
	// 也可以使用sdcodegen.HasSuffix(...).Not() 来指定加载除了.sh和.go文件之外的文件
	// 类似的函数有HasPrefix,HasSuffix,StringContains,FilenameMatches
	err := g.Generate("</path/to/dir>", sdcodegen.HasSuffix(".sh", ".go"))
	if err != nil {
		panic(err)
	}
}
