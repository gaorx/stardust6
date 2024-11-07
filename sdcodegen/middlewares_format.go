package sdcodegen

import (
	"strings"
)

// SetNewline 设置换行符的中间件
func SetNewline(nl string) Middleware {
	return func(c *Context, next Handler) {
		c.SetNewline(nl)
		next(c)
	}
}

// ResetNewline 重置换行符的中间件
func ResetNewline(nl string) Middleware {
	return func(c *Context, next Handler) {
		next(c)
		bufferedText := c.BufferedText()
		oldNL := guessNL(bufferedText)
		if nl != oldNL {
			newText := strings.ReplaceAll(bufferedText, oldNL, nl)
			c.Clear().Print(newText)
		}
	}
}

// FinalNewline 向文件尾部添加换行符的中间件
func FinalNewline() Middleware {
	return func(c *Context, next Handler) {
		next(c)
		if !strings.HasSuffix(c.BufferedText(), c.nl) {
			c.Newl()
		}
	}
}
