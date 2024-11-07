package sdcodegen

// AddHeader 向文件头部添加文本的中间件
func AddHeader(s string) Middleware {
	return func(c *Context, next Handler) {
		c.Print(s)
		next(c)
	}
}

// AddFooter 向文件尾部添加文本的中间件
func AddFooter(s string) Middleware {
	return func(c *Context, next Handler) {
		next(c)
		c.Print(s)
	}
}
