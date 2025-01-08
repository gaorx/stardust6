package sdblueprint

import "github.com/samber/lo"

type Names interface {
	Get(lang string, langAliases ...string) string
	Langs() []string
	Go() string
	Sql() string
	Json() string
	Js() string
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

func (ns *names) Langs() []string {
	return makeLangs(lo.Keys(ns.names), lo.Keys(ns.defaults))
}

func (ns *names) Go() string {
	return ns.Get("go")
}

func (ns *names) Sql() string {
	return ns.Get("sql")
}

func (ns *names) Json() string {
	return ns.Get("json", "js")
}

func (ns *names) Js() string {
	return ns.Get("js", "json")
}

func (ns *names) add(names map[string]string) {
	if len(names) <= 0 {
		return
	}
	for lang, name := range names {
		if lang != "" && name != "" {
			if ns.names == nil {
				ns.names = map[string]string{}
			}
			ns.names[lang] = name
		}
	}
}

func (ns *names) mergeOther(other *names) {
	if len(other.names) <= 0 {
		return
	}
	if ns.names == nil {
		ns.names = map[string]string{}
	}
	for lang, otherName := range other.names {
		if _, ok := ns.names[lang]; !ok {
			ns.names[lang] = otherName
		}
	}
}

func mergeNamers(m, other map[string]Namer) {
	for lang, namer := range other {
		if _, ok := m[lang]; !ok {
			m[lang] = namer
		}
	}
}
