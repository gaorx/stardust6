package sderr

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// PrintOptions 用于控制Print函数的输出
type PrintOptions struct {
	Unwrap         bool // 是否展开多层wrap
	Stack          bool
	FrameFormatter func(f Frame) string
}

// Print 打印一个error到标准输出
func Print(err error, opts *PrintOptions) {
	_ = Fprint(os.Stdout, err, opts)
}

// Sprint 返回一个error的字符串
func Sprint(err error, opts *PrintOptions) string {
	var b strings.Builder
	_ = Fprint(&b, err, opts)
	return b.String()
}

// Fprint 打印一个error到一个输出流
func Fprint(w io.Writer, err error, opts *PrintOptions) error {
	if err == nil {
		return nil
	}

	var errPrint error
	opts = ensurePtr(opts)
	if opts.FrameFormatter == nil {
		opts.FrameFormatter = func(f Frame) string {
			return "  - " + f.String()
		}
	}
	printStack := func(stack *Stack) error {
		frames := stack.Frames()
		for _, frame := range frames {
			_, errPrint0 := fmt.Fprintln(w, opts.FrameFormatter(frame))
			if errPrint0 != nil {
				return Wrap(errPrint0)
			}
		}
		return nil
	}

	if !opts.Unwrap {
		_, errPrint = fmt.Fprintln(w, err.Error())
		if errPrint != nil {
			return Wrap(errPrint)
		}
		if opts.Stack {
			errPrint = printStack(RootStack(err))
			return Wrap(errPrint)
		}
		return nil
	} else {
		unwrappedErrs := UnwrapNested(err)
		for i := 0; i < len(unwrappedErrs); i++ {
			unwrappedErr := unwrappedErrs[i]
			prefix := ""
			if i < len(unwrappedErrs)-1 {
				prefix = "WRAP: "
			} else {
				prefix = "ROOT: "
			}
			if e, ok := Probe(unwrappedErr); ok {
				_, errPrint = fmt.Fprintln(w, prefix+makeMsgWithAttrs(e.Msg(), e.OwnAttrs()))
			} else {
				_, errPrint = fmt.Fprintln(w, prefix+(unwrappedErr.Error()))
			}
			if errPrint != nil {
				return Wrap(errPrint)
			}
			if opts.Stack {
				if e, ok := Probe(unwrappedErr); ok {
					errPrint = printStack(e.stack)
					if errPrint != nil {
						return Wrap(errPrint)
					}
				}
			}
		}
		return nil
	}
}

func makeMsgWithAttrs(msg string, attrs map[string]any) string {
	var b strings.Builder
	b.WriteString(msg)
	if len(attrs) > 0 {
		if msg != "" {
			b.WriteString(" ")
		}
		b.WriteString("[")
		i := 0
		for k, v := range attrs {
			if i > 0 {
				b.WriteString(" ")
			}
			i += 1
			b.WriteString(k)
			b.WriteString("=")
			b.WriteString(fmt.Sprintf("%v", v))
		}
		b.WriteString("]")
	}
	return b.String()
}
