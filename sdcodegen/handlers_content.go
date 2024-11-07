package sdcodegen

// Bytes 向文件中写入一组字节内容
func Bytes(d []byte) Handler {
	return func(c *Context) {
		c.WriteBytes(d)
	}
}

// Text 向文件中写入一组文本内容
func Text(s string) Handler {
	return func(c *Context) {
		c.WriteText(s)
	}
}

// Line 向文件中写入一行文本内容
func Line(s string) Handler {
	return func(c *Context) {
		c.WriteText(s).Newl()
	}
}
