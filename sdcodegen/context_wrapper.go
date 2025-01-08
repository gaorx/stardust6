package sdcodegen

// ContextWrapper 用于包装Context，使得Context的方法可以链式调用
type ContextWrapper[R any] struct {
	c *Context
	r R
}

// MakeContextWrapper 创建一个ContextWrapper
func MakeContextWrapper[R any](c *Context, r R) ContextWrapper[R] {
	return ContextWrapper[R]{c: c, r: r}
}

// TabStyle 获取缩进样式，可能是\t或者几个空格
func (cw ContextWrapper[R]) TabStyle() string {
	return cw.c.TabStyle()
}

// NewlineStyle 获取换行样式，可能是\n或者\r\n
func (cw ContextWrapper[R]) NewlineStyle() string {
	return cw.c.NewlineStyle()
}

// Name 获取生成文件名
func (cw ContextWrapper[R]) Name() string {
	return cw.c.Name()
}

// IsExecutable 判断要生成的文件是否有可执行属性，生成shell脚本时会用到
func (cw ContextWrapper[R]) IsExecutable() bool {
	return cw.c.IsExecutable()
}

// IsDiscarded 判断是否丢弃生成的文件
func (cw ContextWrapper[R]) IsDiscarded() bool {
	return cw.c.IsDiscarded()
}

// Current 获取当前已经存在的文件
func (cw ContextWrapper[R]) Current() *File {
	return cw.c.Current()
}

// BufferedBytes 获取当前已经写入的数据
func (cw ContextWrapper[R]) BufferedBytes() []byte {
	return cw.c.BufferedBytes()
}

// BufferedText 获取当前已经写入的文本
func (cw ContextWrapper[R]) BufferedText() string {
	return cw.c.BufferedText()
}

// SetTab 设置文本缩进样式，例如\t或者几个空格
func (cw ContextWrapper[R]) SetTab(s string) R {
	cw.c.SetTab(s)
	return cw.r
}

// SetNewline 设置文本换行样式，例如\n或者\r\n
func (cw ContextWrapper[R]) SetNewline(s string) R {
	cw.c.SetNewline(s)
	return cw.r
}

// SetExecutable 设置生成文件是否有可执行属性
func (cw ContextWrapper[R]) SetExecutable(b bool) R {
	cw.c.SetExecutable(b)
	return cw.r
}

// SetDiscarded 设置是否丢弃生成的文件
func (cw ContextWrapper[R]) SetDiscarded(b bool) R {
	cw.c.SetDiscarded(b)
	return cw.r
}

// Discard 丢弃生成的文件，相当于设置SetDiscarded(true)
func (cw ContextWrapper[R]) Discard() R {
	cw.c.Discard()
	return cw.r
}

// SetLogger 设置日志
func (cw ContextWrapper[R]) SetLogger(logger Logger) R {
	cw.c.SetLogger(logger)
	return cw.r
}

// SetError 设置运行时产生的错误，一旦产生错误，生成会失败并返回此错误
func (cw ContextWrapper[R]) SetError(err error) R {
	cw.c.SetError(err)
	return cw.r
}

// Abort 中断生成，这之后的代码将不再执行，相当于panic
func (cw ContextWrapper[R]) Abort() {
	cw.c.Abort()
}

// AbortIf 如果条件为真，中断生成
func (cw ContextWrapper[R]) AbortIf(b bool) {
	cw.c.AbortIf(b)
}

// DiscardAndAbortIfExists 如果当前文件已经存在，丢弃此文件并中断生成
func (cw ContextWrapper[R]) DiscardAndAbortIfExists() {
	cw.c.DiscardAndAbortIfExists()
}

// Log 打印日志
func (cw ContextWrapper[R]) Log(msg string) {
	cw.c.Log(msg)
}

// Logf 格式化打印日志
func (cw ContextWrapper[R]) Logf(format string, a ...any) {
	cw.c.Logf(format, a...)
}

// Apply 执行一个函数
func (cw ContextWrapper[R]) Apply(f func()) R {
	cw.c.Apply(f)
	return cw.r
}

// WriteBytes 写入字节数据
func (cw ContextWrapper[R]) WriteBytes(data []byte) R {
	cw.c.WriteBytes(data)
	return cw.r
}

// WriteText 写入文本
func (cw ContextWrapper[R]) WriteText(text string) R {
	cw.c.WriteText(text)
	return cw.r
}

// Clear 清空已经写入的数据
func (cw ContextWrapper[R]) Clear() R {
	cw.c.Clear()
	return cw.r
}

// Tab 打印一个缩进
func (cw ContextWrapper[R]) Tab() R {
	cw.c.Tab()
	return cw.r
}

// TabX 打印n个缩进
func (cw ContextWrapper[R]) TabX(n int) R {
	cw.c.TabX(n)
	return cw.r
}

// Newl 打印一个换行
func (cw ContextWrapper[R]) Newl() R {
	cw.c.Newl()
	return cw.r
}

// NewlX 打印n个换行
func (cw ContextWrapper[R]) NewlX(n int) R {
	cw.c.NewlX(n)
	return cw.r
}

// Print 打印任意数据，这些数据是连接起来的，没有任何分隔符
func (cw ContextWrapper[R]) Print(a ...any) R {
	cw.c.Print(a...)
	return cw.r
}

// Printf 格式化打印
func (cw ContextWrapper[R]) Printf(format string, a ...any) R {
	cw.c.Printf(format, a...)
	return cw.r
}

// Line 打印一行，相当于Print(a...).Newl()
func (cw ContextWrapper[R]) Line(a ...any) R {
	cw.c.Line(a...)
	return cw.r
}

// Linef 格式化打印一行，相当于Printf(format, a...).Newl()
func (cw ContextWrapper[R]) Linef(format string, a ...any) R {
	cw.c.Linef(format, a...)
	return cw.r
}

// If 如果条件为真，打印数据
func (cw ContextWrapper[R]) If(b bool, v any) R {
	cw.c.If(b, v)
	return cw.r
}

// Iff 如果条件为真，打印格式化数据
func (cw ContextWrapper[R]) Iff(b bool, format string, a ...any) R {
	cw.c.Iff(b, format, a...)
	return cw.r
}

// IfElse 如果条件为真，打印v1，否则打印v2
func (cw ContextWrapper[R]) IfElse(b bool, v1, v2 any) R {
	cw.c.IfElse(b, v1, v2)
	return cw.r
}

// ForEach 遍历数组或切片，对每个元素执行action，可以用来打印切片中的每个元素
func (cw ContextWrapper[R]) ForEach(a any, action func(elem any, i, n int)) R {
	cw.c.ForEach(a, action)
	return cw.r
}

// Placeholder 插入一个占位符，占位符是一个特殊的标记，可以在生成文本时替换成其他内容
func (cw ContextWrapper[R]) Placeholder(name string) R {
	cw.c.Placeholder(name)
	return cw.r
}

// ExpandPlaceholder 定义了某个占位符是如何被展开的，被展开的内容将替换到占位符的位置
func (cw ContextWrapper[R]) ExpandPlaceholder(name string, expand func()) R {
	cw.c.ExpandPlaceholder(name, expand)
	return cw.r
}

// GenerateBytes 通过Handler生成另一个数据
func (cw ContextWrapper[R]) GenerateBytes(g Handler, middlewares ...Middleware) []byte {
	return cw.c.GenerateBytes(g, middlewares...)
}

// GenerateText 通过Handler生成另一个文本
func (cw ContextWrapper[R]) GenerateText(g Handler, middlewares ...Middleware) string {
	return cw.c.GenerateText(g, middlewares...)
}

// CURL 从一个URL获取文本
func (cw ContextWrapper[R]) CURL(url string) string {
	return cw.c.CURL(url)
}
