package sdblueprint

type refsType references

var _ Type = refsType{}

func newRefsType(base Type, refs map[string]string) Type {
	if len(refs) <= 0 {
		return base
	}
	return refsType{base: base, refs: refs}
}

func (r refsType) Kind() TypeKind {
	return r.base.Kind()
}

func (r refsType) Refs() Refs {
	rr := references(r)
	return &rr
}

func (r refsType) WithRefs(langRefs map[string]string) Type {
	return newRefsType(r, langRefs)
}

func (r refsType) WithRef(lang, ref string) Type {
	return newRefsType(r, map[string]string{lang: ref})
}

func (r refsType) Schema() Schema {
	return r.base.Schema()
}

func (r refsType) Elem() Type {
	return r.base.Elem()
}

func (r refsType) MakeArray() Type {
	return arrayType{elem: r}
}
