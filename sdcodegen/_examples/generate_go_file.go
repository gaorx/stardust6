package main

import (
	"fmt"
	"github.com/gaorx/stardust6/sdcodegen"
)

func main() {
	// 下面这个例子生成一个GO文件
	f, err := sdcodegen.Generate("src/main.go", nil, func(c *sdcodegen.Context) {
		c.Line("package main")
		c.Newl()
		c.Line("import \"fmt\"")
		c.Newl()
		c.Line("func main() {")
		c.Tab().Line("\tfmt.Println(\"Hello, world!\")")
		c.Line("}")
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(f.StringText())
}
