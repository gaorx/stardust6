package sderr

import (
	"fmt"
	"maps"
)

type Builder Error

func With(k string, v any) Builder {
	return newBuilder().With(k, v)
}

func WithSome(attrs map[string]any) Builder {
	return newBuilder().WithSome(attrs)
}

func (b Builder) With(k string, v any) Builder {
	if k != "" {
		if v1, ok := evalAttrValue(v); ok {
			b.attrs[k] = v1
		}
	}
	return b
}

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

func (b Builder) Newf(msg string, a ...any) error {
	return b.Wrap(fmt.Errorf(msg, a...))
}

func (b Builder) Wrap(err error) error {
	return b.Wrapf(err, "")
}

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

func (b Builder) Recover(f func()) (err error) {
	return b.Recoverf(f, "")
}

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
