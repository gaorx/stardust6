package main

import (
	"fmt"
	"github.com/gaorx/stardust6/sdcodegen"
)

func main() {
	code, err := sdcodegen.GenerateText(func(c *sdcodegen.Context) {
		// 需要倒入的包
		var importPackages []string

		// 定义名为"import"的占位符如何展开
		c.ExpandPlaceholder("import", func() {
			switch len(importPackages) {
			case 0:
			case 1:
				c.Linef(`import "%s"`, importPackages[0])
			default:
				c.Line("import (")
				c.ForEach(importPackages, func(v any, i, n int) {
					c.Tab().Linef(`"%s"`, v.(string))
				})
				c.Line(")")
			}
		})

		c.Line("package main").Newl()
		c.Placeholder("import").Newl() // 在这里插入占位符
		c.Line("func main() {")
		c.Tab().Line(`fmt.Println("hello world")`)
		c.Line("}")

		// 最后面加入需要import的包
		importPackages = append(importPackages, "fmt", "io/fs")
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(code)
}
