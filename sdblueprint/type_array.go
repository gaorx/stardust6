package sdblueprint

type arrayType struct {
	elem Type
}

var _ Type = arrayType{}

func (a arrayType) Kind() TypeKind {
	return KArray
}

func (a arrayType) Refs() Refs {
	m := map[string]string{}
	for _, lang := range a.elem.Refs().Langs() {
		switch lang {
		case "go":
			m["go"] = "[]" + a.elem.Refs().Get("go")
		case "js":
			m["js"] = "[]" + a.elem.Refs().Get("js")
		}
	}
	return &references{base: a, refs: m}
}

func (a arrayType) WithRefs(langRefs map[string]string) Type {
	return newRefsType(a, langRefs)
}

func (a arrayType) WithRef(lang, ref string) Type {
	return newRefsType(a, map[string]string{lang: ref})
}

func (a arrayType) Schema() Schema {
	return nil
}

func (a arrayType) Elem() Type {
	return a.elem
}

func (a arrayType) MakeArray() Type {
	return arrayType{elem: a}
}
