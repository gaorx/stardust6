package sdblueprint

import "github.com/samber/lo"

type Refs interface {
	Get(lang string, langAliases ...string) string
	Langs() []string
	Go() string
	Sql() string
	Json() string
	Js() string
}

type references struct {
	base Type
	refs map[string]string
}

var _ Refs = (*references)(nil)

func (refs *references) Get(lang string, langAliases ...string) string {
	if refs == nil {
		return ""
	}
	langs := append([]string{lang}, langAliases...)
	if r := findLangValue(refs.refs, langs); r != "" {
		return r
	}
	if refs.base != nil {
		if r := refs.base.Refs().Get(lang, langAliases...); r != "" {
			return r
		}
	}
	return ""
}

func (refs *references) Langs() []string {
	var baseLangs []string
	if refs.base != nil {
		baseLangs = refs.base.Refs().Langs()
	}
	return makeLangs(baseLangs, lo.Keys(refs.refs))
}

func (refs *references) Go() string {
	return refs.Get("go")
}

func (refs *references) Sql() string {
	return refs.Get("sql")
}

func (refs *references) Json() string {
	return refs.Get("json", "js")
}

func (refs *references) Js() string {
	return refs.Get("js", "json")
}

func (refs *references) add(langRefs map[string]string) {
	if len(langRefs) <= 0 {
		return
	}
	for lang, v := range langRefs {
		if lang != "" && v != "" {
			if refs.refs == nil {
				refs.refs = make(map[string]string)
			}
			refs.refs[lang] = v
		}
	}
}

func (refs *references) mergeRefs(otherRefs map[string]string) {
	if len(otherRefs) <= 0 {
		return
	}
	for lang, v := range otherRefs {
		if lang != "" && v != "" {
			if refs.refs == nil {
				refs.refs = make(map[string]string)
			}
			if _, ok := refs.refs[lang]; !ok {
				refs.refs[lang] = v
			}
		}
	}
}
