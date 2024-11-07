package main

import (
	"fmt"
	"github.com/gaorx/stardust6/sdcodegen"
	"github.com/gaorx/stardust6/sderr"
	"time"
)

func main() {
	// 生成文件的时候，如果已经有文件存在，通常有3个选择:
	// 1. 覆盖
	// 2. 不覆盖
	// 3. 和已有文件合并
	g := sdcodegen.New()
	g.SetLogger(sdcodegen.Slog(nil))
	g.Add("override.txt", func(c *sdcodegen.Context) {
		// 覆盖
		c.Line("OVERWRITE ", time.Now().Format(time.DateTime))
	})
	g.Add("discard.txt", func(c *sdcodegen.Context) {
		c.DiscardAndAbortIfExists() // 如果文件已经存在，则忽略写入，下面的代码不会执行
		c.Line("DISCARD ", time.Now().Format(time.DateTime))
	})
	g.Add("merge.txt", func(c *sdcodegen.Context) {
		// 合并
		current := c.Current() // 获取当前存在的文件
		if current != nil {
			// 如果文件存在，那么替换掉"// START"和"// END"两行之间的内容，其他部分不变
			text := current.Text() // 获取当前文件内容
			_, err := sdcodegen.ReplaceBetweenTwoLines(
				text,
				sdcodegen.StringContains("// START"),
				sdcodegen.StringContains("// END"),
				func(g *sdcodegen.Context) {
					g.Line("NEW CONTENT ", time.Now().Format(time.DateTime))
				},
			)
			if err != nil {
				panic(sderr.Wrap(err))
			}
		} else {
			// 如果文件不存在，直接写入
			c.Line("HEADER")
			c.Line("// START")
			c.Line("CONTENT ", time.Now().Format(time.DateTime))
			c.Line("// END")
			c.Line("FOOTER")
		}
	})
	fmt.Println("TODO: 接触下行注释，然后修改测试目录，执行生成代码...")
	// g.Generate("</path/to/dir>")
}
