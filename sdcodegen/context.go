package sdcodegen

import (
	"bytes"
	"fmt"
	"github.com/gaorx/stardust6/sderr"
	"github.com/gaorx/stardust6/sdreflect"
	"github.com/samber/lo"
	"reflect"
)

// Context 生成文本时用的上下文
type Context struct {
	tab          string
	nl           string
	name         string
	buff         bytes.Buffer
	executable   bool
	discarded    bool
	current      *File
	err          error
	logger       Logger
	placeholders []*placeholder
	buffering    string
}

type placeholder struct {
	name   string
	stamp  string
	expand func()
	buff   bytes.Buffer
}

// TabStyle 获取缩进样式，可能是\t或者几个空格
func (c *Context) TabStyle() string {
	return c.tab
}

// NewlineStyle 获取换行样式，可能是\n或者\r\n
func (c *Context) NewlineStyle() string {
	return c.nl
}

// Name 获取生成文件名
func (c *Context) Name() string {
	return c.name
}

// IsExecutable 判断要生成的文件是否有可执行属性，生成shell脚本时会用到
func (c *Context) IsExecutable() bool {
	return c.executable
}

// IsDiscarded 判断是否丢弃生成的文件
func (c *Context) IsDiscarded() bool {
	return c.discarded
}

// Current 获取当前已经存在的文件
func (c *Context) Current() *File {
	return c.current
}

// BufferedBytes 获取当前已经写入的数据
func (c *Context) BufferedBytes() []byte {
	return c.buff.Bytes()
}

// BufferedText 获取当前已经写入的文本
func (c *Context) BufferedText() string {
	return c.buff.String()
}

// SetTab 设置文本缩进样式，例如\t或者几个空格
func (c *Context) SetTab(s string) *Context {
	c.tab = s
	return c
}

// SetNewline 设置文本换行样式，例如\n或者\r\n
func (c *Context) SetNewline(s string) *Context {
	c.nl = s
	return c
}

// SetExecutable 设置生成文件是否有可执行属性
func (c *Context) SetExecutable(b bool) *Context {
	c.executable = b
	return c
}

// SetDiscarded 设置是否丢弃生成的文件
func (c *Context) SetDiscarded(b bool) *Context {
	c.discarded = b
	return c
}

// Discard 丢弃生成的文件，相当于设置SetDiscarded(true)
func (c *Context) Discard() *Context {
	c.discarded = true
	return c
}

// SetLogger 设置日志
func (c *Context) SetLogger(logger Logger) *Context {
	if logger == nil {
		logger = NoLog()
	}
	c.logger = logger
	return c
}

// SetError 设置运行时产生的错误，一旦产生错误，生成会失败并返回此错误
func (c *Context) SetError(err error) *Context {
	c.err = err
	return c
}

// Abort 中断生成，这之后的代码将不再执行，相当于panic
func (c *Context) Abort() {
	panic(abortError{})
}

// AbortIf 如果条件为真，中断生成
func (c *Context) AbortIf(b bool) {
	if b {
		c.Abort()
	}
}

// DiscardAndAbortIfExists 如果当前文件已经存在，丢弃此文件并中断生成
func (c *Context) DiscardAndAbortIfExists() {
	if c.current != nil {
		c.Discard().Abort()
	}
}

// Log 打印日志
func (c *Context) Log(msg string) {
	c.logger.Log(msg)
}

// Logf 格式化打印日志
func (c *Context) Logf(format string, a ...any) {
	c.logger.Log(fmt.Sprintf(format, a...))
}

// WriteBytes 写入字节数据
func (c *Context) WriteBytes(data []byte) *Context {
	buff := c.getBuff()
	if buff != nil {
		_, _ = buff.Write(data)
	}
	return c
}

// WriteText 写入文本
func (c *Context) WriteText(text string) *Context {
	return c.WriteBytes([]byte(text))
}

// Clear 清空已经写入的数据
func (c *Context) Clear() *Context {
	c.buff.Reset()
	return c
}

// Tab 打印一个缩进
func (c *Context) Tab() *Context {
	return c.WriteText(c.tab)
}

// TabX 打印n个缩进
func (c *Context) TabX(n int) *Context {
	for i := 0; i < n; i++ {
		c.WriteText(c.tab)
	}
	return c
}

// Newl 打印一个换行
func (c *Context) Newl() *Context {
	return c.WriteText(c.nl)
}

// NewlX 打印n个换行
func (c *Context) NewlX(n int) *Context {
	for i := 0; i < n; i++ {
		c.WriteText(c.nl)
	}
	return c
}

// Print 打印任意数据，这些数据是连接起来的，没有任何分隔符
func (c *Context) Print(a ...any) *Context {
	for _, v := range a {
		switch v1 := v.(type) {
		case nil:
		case string:
			c.WriteText(v1)
		case func():
			if v1 != nil {
				v1()
			}
		default:
			c.WriteText(fmt.Sprintf("%v", v))
		}
	}
	return c
}

// Printf 格式化打印
func (c *Context) Printf(format string, a ...any) *Context {
	return c.WriteText(fmt.Sprintf(format, a...))
}

// Line 打印一行，相当于Print(a...).Newl()
func (c *Context) Line(a ...any) *Context {
	return c.Print(a...).Newl()
}

// Linef 格式化打印一行，相当于Printf(format, a...).Newl()
func (c *Context) Linef(format string, a ...any) *Context {
	return c.Printf(format, a...).Newl()
}

// If 如果条件为真，打印数据
func (c *Context) If(b bool, v any) *Context {
	if b {
		c.Print(v)
	}
	return c
}

// IfElse 如果条件为真，打印v1，否则打印v2
func (c *Context) IfElse(b bool, v1, v2 any) *Context {
	if b {
		c.Print(v1)
	} else {
		c.Print(v2)
	}
	return c
}

// ForEach 遍历数组或切片，对每个元素执行action，可以用来打印切片中的每个元素
func (c *Context) ForEach(a any, action func(elem any, i, n int)) *Context {
	a1 := sdreflect.RootValueOf(a)
	if !sdreflect.IsSliceLikeValue(a1) {
		panic(sderr.Newf("not a slice or array"))
	}
	sdreflect.ForEach(a1, func(elem reflect.Value, i, n int) {
		action(elem.Interface(), i, n)
	})
	return c
}

// Placeholder 插入一个占位符，占位符是一个特殊的标记，可以在生成文本时替换成其他内容
func (c *Context) Placeholder(name string) *Context {
	if name == "" {
		panic(sderr.Newf("placeholder name is empty"))
	}
	p := c.ensurePlaceholder(name)
	return c.WriteText(p.stamp)
}

// ExpandPlaceholder 定义了某个占位符是如何被展开的，被展开的内容将替换到占位符的位置
func (c *Context) ExpandPlaceholder(name string, expand func()) *Context {
	if name == "" {
		panic(sderr.Newf("placeholder name is empty"))
	}
	p := c.ensurePlaceholder(name)
	p.expand = expand
	return c
}

// GenerateBytes 通过Handler生成另一个数据
func (c *Context) GenerateBytes(g Handler, middlewares ...Middleware) []byte {
	data, err := GenerateBytes(g, middlewares...)
	if err != nil {
		panic(sderr.Wrap(err))
	}
	return data
}

// GenerateText 通过Handler生成另一个文本
func (c *Context) GenerateText(g Handler, middlewares ...Middleware) string {
	text, err := GenerateText(g, middlewares...)
	if err != nil {
		panic(sderr.Wrap(err))
	}
	return text
}

func (c *Context) getBuff() *bytes.Buffer {
	if c.buffering == "" {
		return &c.buff
	} else {
		p := c.getPlaceholder(c.buffering)
		if p == nil {
			return nil
		}
		return &p.buff
	}
}

func (c *Context) getPlaceholder(name string) *placeholder {
	for _, p := range c.placeholders {
		if p.name == name {
			return p
		}
	}
	return nil
}

func (c *Context) ensurePlaceholder(name string) *placeholder {
	p := c.getPlaceholder(name)
	if p != nil {
		return p
	}
	stamp := fmt.Sprintf("<<%s:%s>>", name, lo.RandomString(6, lo.AlphanumericCharset))
	p = &placeholder{name: name, stamp: stamp, expand: nil}
	c.placeholders = append(c.placeholders, p)
	return p
}

func (c *Context) expandPlaceholders() {
	placeholders := c.placeholders
	if len(placeholders) <= 0 {
		return
	}
	for _, p := range placeholders {
		func(p *placeholder) {
			c.buffering = p.name
			defer func() { c.buffering = "" }()
			if p.expand != nil {
				p.expand()
			}
			placeholderData, currentData := p.buff.Bytes(), c.BufferedBytes()
			newData := bytes.ReplaceAll(currentData, []byte(p.stamp), placeholderData)
			c.buff.Reset()
			c.buff.Write(newData)
		}(p)
	}
}
