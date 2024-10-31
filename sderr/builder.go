package sderr

import (
	"fmt"
	"maps"
)

type Builder Error

// With 创建一个error builder，带有一个Attr
func With(k string, v any) Builder {
	return newBuilder().With(k, v)
}

// WithSome 创建一个error builder，带有多个Attr
func WithSome(attrs map[string]any) Builder {
	return newBuilder().WithSome(attrs)
}

// With 附加一个Attr到Builder
func (b Builder) With(k string, v any) Builder {
	if k != "" {
		if v1, ok := evalAttrValue(v); ok {
			b.attrs[k] = v1
		}
	}
	return b
}

// WithSome 附加多个Attr到Builder
func (b Builder) WithSome(attrs map[string]any) Builder {
	for k, v := range attrs {
		if k != "" {
			if v1, ok := evalAttrValue(v); ok {
				b.attrs[k] = v1
			}
		}
	}
	return b
}

// Newf 构建一个带有stacktrace的error，并携带一个可读消息
func (b Builder) Newf(msg string, a ...any) error {
	return b.Wrap(fmt.Errorf(msg, a...))
}

// Wrap 构建一个带有stacktrace的error，用于wrap一个现存的error
func (b Builder) Wrap(err error) error {
	return b.Wrapf(err, "")
}

// Wrapf 构建一个带有stacktrace的error，用于wrap一个现存的error，并附加一个可读消息
func (b Builder) Wrapf(err error, msg string, a ...any) error {
	if err == nil {
		return nil
	}
	b2 := newBuilder()
	b2.err = err
	b2.msg = fmtMsg(msg, a...)
	b2.attrs = copyAttrs(b.attrs)
	b2.stack = newStacktrace()
	return (*Error)(&b2)
}

// Recover 构建一个带有stacktrace的error，并携带一个函数的panic信息
func (b Builder) Recover(f func()) (err error) {
	return b.Recoverf(f, "")
}

// Recoverf 构建一个带有stacktrace的error，并携带一个函数的panic信息和一个可读消息
func (b Builder) Recoverf(f func(), msg string, a ...any) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				err = b.Wrap(e)
			} else {
				err = b.Newf("%v", r)
			}
		}
	}()

	f()
	return
}

func newBuilder() Builder {
	return Builder{
		attrs: map[string]any{},
	}
}

func copyAttrs(src map[string]any) map[string]any {
	dst := map[string]any{}
	maps.Copy(dst, src)
	return dst
}

func fmtMsg(msg string, a ...any) string {
	if len(a) <= 0 {
		return msg
	}
	return fmt.Errorf(msg, a...).Error()
}
