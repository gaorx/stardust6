package sdblueprint

import (
	"github.com/samber/lo"
	"slices"
)

type Anns interface {
	Get(lang string, langAliases ...string) []string
	Langs() []string
	Go() []string
	Js() []string
}

type annotations struct {
	anns map[string][]string
}

var _ Anns = (*annotations)(nil)

func (anns *annotations) Get(lang string, langAliases ...string) []string {
	langs := append([]string{lang}, langAliases...)
	r := findLangValue(anns.anns, langs)
	if r != nil {
		return r
	}
	return nil
}

func (anns *annotations) Langs() []string {
	return lo.Keys(anns.anns)
}

func (anns *annotations) Go() []string {
	return anns.Get("go")
}

func (anns *annotations) Js() []string {
	return anns.Get("js")
}

func (anns *annotations) add(lang string, ann []string) {
	if lang == "" || len(ann) <= 0 {
		return
	}
	if anns.anns == nil {
		anns.anns = map[string][]string{}
	}
	anns.anns[lang] = ann
}

func (anns *annotations) mergeOther(other *annotations) {
	if len(other.anns) <= 0 {
		return
	}
	if anns.anns == nil {
		anns.anns = map[string][]string{}
	}
	for lang, otherAnn := range other.anns {
		if _, ok := anns.anns[lang]; !ok {
			anns.anns[lang] = slices.Clone(otherAnn)
		}
	}
}
