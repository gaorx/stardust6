package sdgogen

import (
	"github.com/samber/lo"
	"strings"
)

type Param struct {
	Name string
	Type string
}

type Params []Param

func P(name, typ string) Param {
	return Param{Name: name, Type: typ}
}

func T(typ string) Param {
	return Param{Type: typ}
}

func (p Param) String() string {
	if p.Type != "" && p.Name != "" {
		return p.Name + " " + p.Type
	} else if p.Name != "" {
		return p.Name
	} else if p.Type != "" {
		return p.Type
	} else {
		return ""
	}
}

func (ps Params) String() string {
	if len(ps) <= 0 {
		return ""
	}
	ss := lo.Map(ps, func(p Param, _ int) string {
		return p.String()
	})
	ss = lo.Filter(ss, func(s string, _ int) bool {
		return s != ""
	})
	return strings.Join(ss, ", ")
}

func (ps Params) StringReturns() string {
	ss := lo.Map(ps, func(p Param, _ int) string {
		return p.String()
	})
	ss = lo.Filter(ss, func(s string, _ int) bool {
		return s != ""
	})
	switch len(ss) {
	case 0:
		return ""
	case 1:
		return " " + ss[0]
	default:
		return " (" + strings.Join(ss, ", ") + ")"
	}
}

func (ps Params) WithErr() Params {
	return append(ps, P("", "error"))
}

func (ps Params) WithNamedErr(name string) Params {
	return append(ps, P(name, "error"))
}
