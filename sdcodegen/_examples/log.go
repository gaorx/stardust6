package main

import (
	"fmt"
	"github.com/gaorx/stardust6/sdcodegen"
	"github.com/gaorx/stardust6/sdslog"
	"log/slog"
)

func main() {
	sdslog.SetDefault([]sdslog.Handler{
		sdslog.TextFile(slog.LevelDebug, "stdout", true),
	}, nil)
	// 可以在生成过程中打印日志，但必须配备SetLogger中间件
	// 在generator中，使用(*Generator).SetLogger()则会自动配置日志，无需再手工配置
	s, err := sdcodegen.GenerateText(func(c *sdcodegen.Context) {
		c.Log("START...")
		c.Line("HELLO WORLD")
		c.Log("END.")
	}, sdcodegen.SetLogger(sdcodegen.Slog(nil)))
	if err != nil {
		panic(err)
	}
	fmt.Println("==============")
	fmt.Println(s)
}
