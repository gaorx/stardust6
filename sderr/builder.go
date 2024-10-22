package sderr

import (
	"fmt"
	"maps"
)

type Builder Error

func With(k string, v any) Builder {
	return newBuilder().With(k, v)
}

func WithAttrs(attrs map[string]any) Builder {
	return newBuilder().WithAttrs(attrs)
}

func (b Builder) With(k string, v any) Builder {
	b.attrs[k] = v
	return b
}

func (b Builder) WithAttrs(attrs map[string]any) Builder {
	for k, v := range attrs {
		b.attrs[k] = v
	}
	return b
}

func (b Builder) New(msg string) error {
	return Wrap(nil, msg)
}

func (b Builder) Newf(msg string, a ...any) error {
	return b.New(fmt.Sprintf(msg, a...))
}

func (b Builder) Wrap(err error, msg string) error {
	b2 := newBuilder()
	b2.err = err
	b2.msg = msg
	b2.attrs = copyAttrs(b.attrs)
	b2.stack = newStacktrace()
	return Error(b2)
}

func (b Builder) Wrapf(err error, msg string, a ...any) error {
	return b.Wrap(err, fmt.Sprintf(msg, a...))
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
