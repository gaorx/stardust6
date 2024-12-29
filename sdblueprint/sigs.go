package sdblueprint

import (
	"github.com/gaorx/stardust6/sdstrings"
	"github.com/samber/lo"
	"strings"
)

type Sig struct {
	Pkg string
	Typ string
}

func NewSig(pkg, typ string) Sig {
	return Sig{pkg, typ}
}

func SigOf(s string) Sig {
	if s == "" {
		return Sig{}
	}
	if !strings.Contains(s, "#") {
		return NewSig("", s)
	}
	pkg, typ := sdstrings.Split2s(s, "#")
	return NewSig(pkg, typ)
}

type Sigs map[string]Sig

func (lt Sigs) Get(lang string, langAliases ...string) Sig {
	langs := append([]string{lang}, langAliases...)
	r := findLangValue(lt, langs)
	if !lo.IsEmpty(r) {
		return r
	}
	var zero Sig
	return zero
}
