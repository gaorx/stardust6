package sdprint

import (
	"fmt"
	"github.com/gaorx/stardust6/sdjson"
	"github.com/k0kubun/pp/v3"
	"io"
	"os"
)

// Printer 打印器
type Printer struct {
	W io.Writer
}

// Out 创建一个打印器
func Out(w io.Writer) Printer {
	return Printer{W: w}
}

// Stdout 创建一个标准输出打印器
func Stdout() Printer {
	return Out(os.Stdout)
}

// Stderr 创建一个标准错误打印器
func Stderr() Printer {
	return Out(os.Stderr)
}

// Print 打印一组数据，中间不插入空格
func (p Printer) Print(a ...any) Printer {
	_, _ = fmt.Fprint(p.W, a...)
	return p
}

// Println 打印一组数据并换行，中间插入空格
func (p Printer) Println(a ...any) Printer {
	_, _ = fmt.Fprintln(p.W, a...)
	return p
}

// Printf 格式化打印
func (p Printer) Printf(format string, a ...any) Printer {
	_, _ = fmt.Fprintf(p.W, format, a...)
	return p
}

// Ln 打印换行
func (p Printer) Ln() Printer {
	_, _ = fmt.Fprintln(p.W)
	return p
}

// Json 打印json
func (p Printer) Json(v any) Printer {
	return p.Print(sdjson.MarshalPretty(v))
}

// Jsonln 打印json并换行
func (p Printer) Jsonln(v any) Printer {
	return p.Println(sdjson.MarshalPretty(v))
}

// PlainJson 打印json，不带格式
func (p Printer) PlainJson(v any) Printer {
	return p.Print(sdjson.MarshalStringOr(v, ""))
}

// PlainJsonln 打印json并换行，不带格式
func (p Printer) PlainJsonln(v any) Printer {
	return p.Println(sdjson.MarshalStringOr(v, ""))
}

// Pretty 打印pretty格式
func (p Printer) Pretty(v any) Printer {
	return p.Print(pp.Sprint(v))
}

// Prettyln 打印pretty格式并换行
func (p Printer) Prettyln(v any) Printer {
	return p.Println(pp.Sprint(v))
}
