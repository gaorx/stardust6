package sdcodegen

// Formatter 格式化代码
type Formatter func(code string) (string, error)

// FormatterSelector 选择器
type FormatterSelector func(c *Context) Formatter

// AsMiddleware 使用此formatter去作为中间件格式化生成的代码
func (f Formatter) AsMiddleware() Middleware {
	return FormatSource(func(c *Context) Formatter {
		return f
	})
}

// SelectFormatter 执行此选择器，选择一个formatter
func (sel FormatterSelector) SelectFormatter(c *Context) Formatter {
	if sel == nil {
		return nil
	}
	return sel(c)
}
