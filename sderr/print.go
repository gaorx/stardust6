package sderr

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type PrintOptions struct {
	Unwrap         bool
	Stack          bool
	FrameFormatter func(f Frame) string
}

func Print(err error, opts *PrintOptions) {
	_ = Fprint(os.Stdout, err, opts)
}

func Sprint(err error, opts *PrintOptions) string {
	var b strings.Builder
	_ = Fprint(&b, err, opts)
	return b.String()
}

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
		if e, ok := Probe(err); ok {
			_, errPrint = fmt.Fprintln(w, makeMsgWithAttrs("", e.Error(), Attrs(e)))
		} else {
			_, errPrint = fmt.Fprintln(w, quote(err.Error()))
		}
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
				_, errPrint = fmt.Fprintln(w, makeMsgWithAttrs(prefix, e.Message(), e.OwnAttrs()))
			} else {
				_, errPrint = fmt.Fprintln(w, prefix+quote(unwrappedErr.Error()))
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

func makeMsgWithAttrs(prefix, msg string, attrs map[string]any) string {
	var b strings.Builder
	b.WriteString(prefix)
	b.WriteString(quote(msg))
	if len(attrs) > 0 {
		for k, v := range attrs {
			b.WriteString(" [")
			b.WriteString(k)
			b.WriteString("=")
			b.WriteString(quote(fmt.Sprintf("%v", v)))
			b.WriteString("]")
		}
	}
	return b.String()
}
