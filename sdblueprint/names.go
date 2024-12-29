package sdblueprint

import (
	"strings"
)

type Names interface {
	Get(lang string, langAliases ...string) string
}

var _ Names = (*names)(nil)

type names struct {
	id       string
	names    map[string]string
	defaults map[string]Namer
}

func (ns *names) Get(lang string, langAliases ...string) string {
	langs := append([]string{lang}, langAliases...)
	r := findLangValue(ns.names, langs)
	if r != "" {
		return r
	}
	def := findLangValue(ns.defaults, langs)
	if def != nil {
		return def(ns.id)
	}
	return ns.id
}

func (ns *names) add(langAndNames []string) {
	if len(langAndNames) <= 0 {
		return
	}
	for lang, v := range makeLangMap(langAndNames) {
		if ns.names == nil {
			ns.names = make(map[string]string)
		}
		ns.names[lang] = v
	}
}

func findLangValue[V string | []string | Namer | Sig](m map[string]V, langs []string) V {
	var zero V
	if len(langs) <= 0 {
		return zero
	}
	for _, lang := range langs {
		for k, v := range m {
			if k != "" && lang != "" {
				if strings.ToLower(k) == strings.ToLower(lang) {
					return v
				}
			}
		}
	}
	return zero
}
