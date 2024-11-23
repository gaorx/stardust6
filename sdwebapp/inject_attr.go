package sdwebapp

import (
	"github.com/samber/lo"
	"slices"
)

type Attribute struct {
	K string
	V any
}

type Attrs []Attribute
type AttrMap map[string]any

var (
	_ Injectable = Attribute{}
	_ Injectable = AttrMap{}
	_ Injectable = Attrs(nil)
)

func Attr(k string, v any) Attribute {
	return Attribute{K: k, V: v}
}

func (a Attribute) Injections() Attrs {
	return []Attribute{{K: a.K, V: a.V}}
}

func (attrs Attrs) Injections() Attrs {
	return slices.Clone(attrs)
}

func (m AttrMap) Injections() Attrs {
	return lo.MapToSlice(m, func(k string, v any) Attribute {
		return Attr(k, v)
	})
}
